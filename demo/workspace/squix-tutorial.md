---
title: "Learning Squix"
author: "Squix Demo"
---

# Welcome to Squix! 🐿️

**Squix's SQL Stash**

A query manager for your databases

<!-- end_slide -->

# What is Squix?

Squix is a powerful **query manager** that helps you:

- 🔗 **Connect** to multiple databases
- 🔍 **Explore** database schemas
- 💾 **Save** queries for quick access
- 🚀 **Run** queries with simple commands
- 📋 **Organize** queries into runbooks
- 🌳 **Visualize** table relationships

Perfect for:
- Data analysis
- Database exploration
- Query documentation
- Team collaboration

<!-- end_slide -->

# Our Dataset: Squirrels! 🐿️

We have a SQLite database with squirrel tracking data:

**squirrels** table: 10 squirrels
- Name, species, age, favorite food, location

**sightings** table: 20 observations
- Date, behavior, notes

Let's explore it with Squix!

<!-- end_slide -->

# Step 1: Initialize Connection

First, tell Squix about our database:

```bash +exec
squix init --name squirrels --type sqlite --conn "/workspace/squirrels.db"
```
This creates a connection named "squirrels" to our SQLite database.

<!-- end_slide -->

# Step 2: Discover Tables

See what tables are available:

```bash +exec +pty
squix tables
```

Lists all tables in the current database!

<!-- end_slide -->

# Step 3: Quick Table Preview

Query any table instantly:

```bash +exec +pty
squix tables squirrels
```

No need to type SELECT * - Squix does it for you!

<!-- end_slide -->

# Step 4: Explore Schema

Get an overview of your database:

```bash +exec +pty
squix explore
```

Shows all tables and views in a clean format.

<!-- end_slide -->

# Step 5: Quick Query

Or run raw SQL when needed:

```bash +exec +pty
squix query "SELECT * FROM squirrels LIMIT 5"
```

See the first 5 squirrels!

<!-- end_slide -->

# Step 6: List Connections

Managing multiple databases? Check your connections:

```bash +exec +pty
squix list connections
```

Shows all configured database connections.

<!-- end_slide -->

# Step 7: Connection Info

Get details about the active connection:

```bash +exec +pty
squix info
```

Shows tables, views, and metadata.

<!-- end_slide -->

# Step 8: Create a Saved Query

Save a query for reuse:

```bash +exec
squix add "all-squirrels" "SELECT * FROM squirrels"
```

Now you can run `squix run all-squirrels` anytime!

<!-- end_slide -->

# Step 9: Run a Saved Query

Execute your saved query:

```bash +exec +pty
squix run "all-squirrels"
```

Much faster than typing the full SQL every time!

<!-- end_slide -->

# Step 10: Edit Saved Queries

Need to modify a query? Open in your editor:

```bash +exec
squix edit queries
```

Opens the queries file in vim for easy editing.

<!-- end_slide -->

# Step 11: Filter Query

Let's save a filtered query:

```bash +exec
squix add "young-squirrels" "SELECT name, age_years FROM squirrels WHERE age_years < 3"
```

```bash +exec +pty
squix run "young-squirrels"
```

<!-- end_slide -->

# Step 12: Complex JOIN Query

Let's join squirrels with their sightings:

```bash +exec
squix add "recent-sightings" "SELECT s.name, sp.date, sp.behavior FROM squirrels s JOIN sightings sp ON s.id = sp.squirrel_id ORDER BY sp.date DESC"
```

```bash +exec +pty
squix run "recent-sightings"
```

<!-- end_slide -->

# Step 13: Table Relationships

Visualize how tables connect:

```bash +exec +pty
squix explain squirrels --depth 1
```

Shows foreign key relationships in a tree view!

<!-- end_slide -->

# Step 14: Statistics Query

Find the most common behaviors:

```bash +exec
squix add "behavior-stats" "SELECT behavior, COUNT(*) as count FROM sightings GROUP BY behavior ORDER BY count DESC"
```

```bash +exec +pty
squix run "behavior-stats"
```

<!-- end_slide -->

# Step 15: List Saved Queries

See all your saved queries:

```bash +exec
squix list queries
```

Shows all the queries you've created!

<!-- end_slide -->

# Step 16: Favorite Foods

What do squirrels love to eat?

```bash +exec
squix add "favorite-foods" "SELECT favorite_food, COUNT(*) as count FROM squirrels GROUP BY favorite_food"
```

```bash +exec +pty
squix run "favorite-foods"
```

<!-- end_slide -->

# Step 17: Remove Query

Clean up queries you don't need:

```bash +exec
squix remove "young-squirrels"
```

```bash +exec
squix list queries
```

<!-- end_slide -->

# Summary: Key Commands

| Command | Description |
|---------|-------------|
| `squix init` | Create database connection |
| `squix switch` | Switch active connection |
| `squix tables` | List all tables |
| `squix tables <name>` | Quick SELECT * FROM table |
| `squix explore` | Explore database schema |
| `squix explain` | Show table relationships |
| `squix add <name> <sql>` | Save a query |
| `squix run <name>` | Execute saved query |
| `squix edit queries` | Edit saved queries |
| `squix remove <name>` | Delete a query |
| `squix list queries` | List all saved queries |
| `squix list connections` | List all connections |
| `squix info` | Show connection details |

<!-- end_slide -->

# Next Steps

🔧 **More Features:**
- Query parameters and templates
- Export results to CSV/JSON
- Multiple database types (Postgres, MySQL, SQL Server, Oracle)
- Runbooks for organizing queries

📚 **Learn More:**
```bash +exec
squix help
```

<!-- end_slide -->

# Try It Yourself!

Now it's your turn to explore:
```bash +exec
# List all squirrels in Central Park
squix add "central-park" "SELECT * FROM squirrels WHERE park_location LIKE '%Central Park%'"

# Find foraging behaviors
squix add "foraging" "SELECT s.name, sp.date FROM squirrels s JOIN sightings sp ON s.id = sp.squirrel_id WHERE sp.behavior = 'foraging'"

# Run them!
squix run "central-park"
squix run "foraging"

# Explore the schema
squix explore sightings
```

<!-- end_slide -->

# Questions?

🐿️ Happy querying with Squix!

```bash +exec
squix help
```

Type `demo` anytime to restart this presentation!

<!-- end_slide -->
