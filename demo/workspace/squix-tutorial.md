---
theme:
  name: catppuccin-frappe
  override:
    footer:
      style: template
      left: '🐿️'
      center: 'squix'
      right: "{current_slide} / {total_slides}"
      height: 5
    palette:
      classes:
        noice:
          foreground: red
---

<!-- no_footer -->
<!-- newlines: 8 -->
Squix's SQL Stash
---
![](squix-mascot.png)

<!-- alignment: center -->
SQL explorer for terminal squirrels

<!-- end_slide -->

# Squirrels hate leaving the terminal

In a nutshell, squix remembers your database connections and queries for you, allowing you to run repetitve and hard-to-remember queries from the comfort of your terminal

<!-- pause -->

To save a db connection, use the init command (we won't actually run it this time because the connection is already setup in the live demo)
```bash
squix init squirrels /home/scrat/tree/squirrels.db
```

Here we connected to a sample sqlite databse and gave it the name "squirrels". You can connect to a number of other database types like postgres, oracle, mysql, sqlserver, etc -> [more info here](https://github.com/eduardofuncao/squix?tab=readme-ov-file#--------database-support)


<!-- pause -->

This will also set this connection up as your *active connection*. Check it out with (run commands from this presentation with CTRL+E): 
```bash +exec +pty
squix status
```

> [!tip]
> use `squix list connections` to list your saved db connections, and `squix use <conn-name>` to set an active connection 

<!-- end_slide -->

# Understand your new home layout

To explore your database, **squix** has some nice QOL features, such as the explore command, which will list all tables and views available in your current connection
```bash +exec
squix explore
```

<!-- pause -->

Or the explain command, which will show foreign key relationships from a target table
```bash +exec
squix explain squirrels
```

see also `squix help tables` for another options to navigate your database schema

<!-- end_slide -->

# Build your SQL stash

Adding saved queries is another way to claim your spot as a certified terminal squirrel. Add them with the `squix add` command

```bash +exec
squix add show_all "select * from squirrels"
```

<!-- pause -->

The queries are saved per connection, and can be listed with `squix list` or filtered with

```bash +exec
squix list show
```

which will use "show" as a search term on the queries names and SQL


<!-- end_slide -->

# Enjoy the fruits of your labour

You can run queries by name or id with `squix run <query-name/id>`

<!-- pause -->
Inside the table results, you can:

## Navigate, copy and export content
- navigate with `hjkl` (or arrow keys) 
- copy cell's contents with `y`, or start a selection with `v` to copy a larger area
- export selections as html, markdown tables, json or csv with `x`

<!-- pause -->

## Update and Delete stuff
All of the following keymaps will open a query to update/delete your date in your $EDITOR of choice (or vim as a fallback). Save and quit (:wq in vim) to apply changes, or just quit (:q in vim) to cancel the operation
- update a single cell with `u`
- delete a row with D (uppercase!)
- rerun or change your query with `e`

<!-- pause -->

### Other vim-like commands also work for naviation
- g/G to jump to first and last rows
- 0/$ to jump to first and last colum
- CTRL+U/CTRL+D to scroll up and down
- q to quit

Run your first squix query (this will suspend the presentation until the squix command finishes; navigate and explore the results as you like)

```bash +exec +acquire_terminal
squix run show_all
```

Alternatively, you can just run squix with inline queries (eg `squix run "select * from squirrels where name like 'A%'"`). This can be useful for quick testing, before adding queries to your stash)




<!-- end_slide -->

# When things don't go as planned

You can use the `edit` and `remove` commands

## Edit
run `squix edit` with no arguments to open all your saved queries from the active connection. You can also use `squix edit <query-name/id>` to target a specific query for editing. This will also use your $EDITOR:

```bash +acquire_terminal
squix edit show_all
```

## Remove
Remove saved queries with `squix remove <query-name>`, or remove all connection data with `squix remove --connection <conn-name>`


<!-- end_slide -->
<!-- newlines: 8 -->

## Thanks for checking out squix 

Now you can end this presentation with `q` or `CTRL+C` and go back to the interactive terminal. Use `squix explore` or `squix explain` to understand the demo database, `add` and `run` queries to try *squix* out. 

If you get stuck, use the help command `squix help` for general info or `squix help <command>` for command specific info

If you liked it so far, please consider giving it a star and installing it from [](https://github.com/eduardofuncao/squix)

<!-- newlines: 8 -->
<!-- alignment: center -->
made with 🐿️ by @eduardofuncao
