package explain

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/eduardofuncao/squix/internal/db"
	"github.com/eduardofuncao/squix/internal/styles"
)

type relationshipType int

const (
	belongsTo relationshipType = iota
	hasMany
	hasOne
	hasManyToMany
)

type relationship struct {
	relType          relationshipType
	column           string
	referencedTable  string
	referencedColumn string
	junctionTable    string // For N:N relationships
}

// BuildTree renders an ASCII tree of foreign key relationships for a table,
// up to maxDepth levels deep. Set verbose to include FK column details.
func BuildTree(conn db.DatabaseConnection, tableName string, maxDepth int, verbose bool) string {
	if maxDepth < 1 {
		maxDepth = 1
	}

	visited := make(map[string]bool)
	return renderNode(conn, tableName, tableRelationships(conn, tableName), maxDepth, 0, visited, true, verbose)
}

// tableRelationships collects the direct relationships of a table.
func tableRelationships(conn db.DatabaseConnection, tableName string) []relationship {
	var relationships []relationship
	seenTables := make(map[string]bool)

	// Get unique constraints for this table (for 1:1 detection)
	uniqueConstraints, _ := conn.GetUniqueConstraints(tableName)
	uniqueMap := make(map[string]bool)
	for _, uc := range uniqueConstraints {
		uniqueMap[uc] = true
	}

	// Get "belongs to" relationships (FKs from this table to other tables)
	belongsToFKs, err := conn.GetForeignKeys(tableName)
	if err == nil {
		for _, fk := range belongsToFKs {
			key := fmt.Sprintf("%s:%s:%s", fk.ReferencedTable, fk.Column, fk.ReferencedColumn)
			if seenTables[key] {
				continue
			}
			seenTables[key] = true

			relType := belongsTo
			if uniqueMap[fk.Column] {
				relType = hasOne
			}

			relationships = append(relationships, relationship{
				relType:          relType,
				column:           fk.Column,
				referencedTable:  fk.ReferencedTable,
				referencedColumn: fk.ReferencedColumn,
			})
		}
	}

	// Get "has many" relationships (FKs from other tables to this table)
	hasManyFKs, err := conn.GetForeignKeysReferencingTable(tableName)
	if err == nil {
		for _, fk := range hasManyFKs {
			// Check if the referencing table is a junction table (N:N detection)
			if isJunctionTable(conn, fk.ReferencedTable) {
				otherTable := getJunctionTableOtherSide(conn, fk.ReferencedTable, tableName)
				if otherTable != "" {
					key := fmt.Sprintf("nn:%s:%s:%s", otherTable, fk.ReferencedTable, tableName)
					if seenTables[key] {
						continue
					}
					seenTables[key] = true

					relationships = append(relationships, relationship{
						relType:         hasManyToMany,
						referencedTable: otherTable,
						junctionTable:   fk.ReferencedTable,
					})
				}
			} else {
				otherTableUnique, _ := conn.GetUniqueConstraints(fk.ReferencedTable)
				isOneToOne := false
				for _, uc := range otherTableUnique {
					if uc == fk.Column {
						isOneToOne = true
						break
					}
				}

				key := fmt.Sprintf("%s:%s:%s", fk.ReferencedTable, fk.Column, fk.ReferencedColumn)
				if seenTables[key] {
					continue
				}
				seenTables[key] = true

				relType := hasMany
				if isOneToOne {
					relType = hasOne
				}

				relationships = append(relationships, relationship{
					relType:          relType,
					column:           fk.Column,
					referencedTable:  fk.ReferencedTable,
					referencedColumn: fk.ReferencedColumn,
				})
			}
		}
	}

	return relationships
}

// isJunctionTable checks if a table is a junction table for N:N relationship
func isJunctionTable(conn db.DatabaseConnection, tableName string) bool {
	fks, err := conn.GetForeignKeys(tableName)
	if err != nil || len(fks) != 2 {
		return false
	}

	// Check if both FKs are part of the primary key (composite PK)
	metadata, err := conn.GetTableMetadata(tableName)
	if err != nil {
		return false
	}

	pkMap := make(map[string]bool)
	for _, pk := range metadata.PrimaryKeys {
		pkMap[pk] = true
	}

	bothInPK := true
	for _, fk := range fks {
		if !pkMap[fk.Column] {
			bothInPK = false
			break
		}
	}

	return bothInPK
}

// getJunctionTableOtherSide returns the other table referenced by a junction table
func getJunctionTableOtherSide(conn db.DatabaseConnection, junctionTable, currentTable string) string {
	fks, err := conn.GetForeignKeys(junctionTable)
	if err != nil || len(fks) != 2 {
		return ""
	}

	for _, fk := range fks {
		if fk.ReferencedTable != currentTable {
			return fk.ReferencedTable
		}
	}

	return ""
}

func renderNode(
	conn db.DatabaseConnection,
	tableName string,
	relationships []relationship,
	maxDepth, currentDepth int,
	visited map[string]bool,
	isRoot bool,
	verbose bool,
) string {
	var builder strings.Builder

	// Render current table with PK info (only at root level)
	if isRoot {
		metadata, _ := conn.GetTableMetadata(tableName)
		if len(metadata.PrimaryKeys) > 0 {
			pks := strings.Join(metadata.PrimaryKeys, ", ")
			builder.WriteString(styles.TableName.Render(tableName))
			builder.WriteString(" ")
			builder.WriteString(styles.PrimaryKeyLabel.Render(fmt.Sprintf("(PK: %s)", pks)))
		} else {
			builder.WriteString(styles.TableName.Render(tableName))
		}
		builder.WriteString("\n")
	}

	// Mark current table as visited to prevent cycles
	visited[tableName] = true

	// Check if we should render relationships at this depth
	if currentDepth >= maxDepth && !isRoot {
		return builder.String()
	}

	for i, rel := range relationships {
		isLast := i == len(relationships)-1
		prefix := "├── "
		if isLast {
			prefix = "└── "
		}

		var relText, cardinality, fkDetails string
		var relStyle lipgloss.Style

		isSelfReference := (rel.referencedTable == tableName)

		switch rel.relType {
		case belongsTo:
			relText = "belongs to"
			cardinality = "[N:1]"
			relStyle = styles.BelongsToStyle
			if verbose {
				fkDetails = fmt.Sprintf("(FK: %s → %s.%s)", rel.column, rel.referencedTable, rel.referencedColumn)
			}
		case hasOne:
			relText = "has one"
			cardinality = "[1:1]"
			relStyle = styles.HasOneStyle
			if verbose {
				fkDetails = fmt.Sprintf("(FK: %s → %s.%s)", rel.column, rel.referencedTable, rel.referencedColumn)
			}
		case hasMany:
			relText = "has many"
			cardinality = "[1:N]"
			relStyle = styles.HasManyStyle
			if verbose {
				fkDetails = fmt.Sprintf("(on: %s ← %s.%s)", rel.referencedColumn, rel.referencedTable, rel.column)
			}
		case hasManyToMany:
			relText = "↔"
			cardinality = "[N:N]"
			relStyle = styles.HasManyToManyStyle
			if verbose {
				fkDetails = fmt.Sprintf("(via %s)", rel.junctionTable)
			}
		}

		builder.WriteString(styles.TreeConnector.Render(prefix))
		if rel.relType == hasManyToMany {
			builder.WriteString(relStyle.Render(relText))
			builder.WriteString(" ")
			builder.WriteString(styles.TableName.Render(rel.referencedTable))
			builder.WriteString(" ")
			builder.WriteString(styles.CardinalityStyle.Render(cardinality))
		} else {
			builder.WriteString(relStyle.Render(fmt.Sprintf("%s →", relText)))
			builder.WriteString(" ")
			builder.WriteString(styles.TableName.Render(rel.referencedTable))
			builder.WriteString(" ")
			builder.WriteString(styles.CardinalityStyle.Render(cardinality))
		}

		if verbose && fkDetails != "" {
			builder.WriteString(" ")
			builder.WriteString(styles.Faint.Render(fkDetails))
		}

		if isSelfReference {
			builder.WriteString(" " + styles.Faint.Render("(self-reference)"))
		}

		builder.WriteString("\n")

		// Don't show children for self-references or N:N relationships (collapsed)
		if isSelfReference || rel.relType == hasManyToMany {
			continue
		}

		// Don't revisit tables
		if visited[rel.referencedTable] {
			continue
		}

		// Recursively render children if within depth limit
		if currentDepth+1 <= maxDepth {
			childRelationships := tableRelationships(conn, rel.referencedTable)
			if len(childRelationships) > 0 {
				childPrefix := "    "
				if !isLast {
					childPrefix = "│   "
				}

				localVisited := make(map[string]bool)
				for k, v := range visited {
					localVisited[k] = v
				}
				localVisited[rel.referencedTable] = true

				childTree := renderNode(conn, rel.referencedTable, childRelationships, maxDepth, currentDepth+1, localVisited, false, verbose)
				lines := strings.Split(childTree, "\n")
				for _, line := range lines {
					if line == "" {
						continue
					}
					builder.WriteString(styles.TreeConnector.Render(childPrefix))
					builder.WriteString(line)
					builder.WriteString("\n")
				}
			}
		}
	}

	return builder.String()
}
