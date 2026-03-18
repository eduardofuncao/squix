# Squix Demo Welcome Message - Colors matching Squix DefaultScheme
PRIMARY='\033[38;5;205m'   # Magenta - headers, titles
SUCCESS='\033[38;5;171m'   # Purple - success messages
ACCENT='\033[38;5;86m'     # Cyan - commands, keywords
MUTED='\033[38;5;238m'     # Gray - subtitles, info
ERROR='\033[38;5;196m'     # Red - errors, warnings
NC='\033[0m'

# Show appropriate banner based on terminal width
if [ "$COLUMNS" -lt 60 ] 2>/dev/null || [ -z "$COLUMNS" ]; then
    # Mobile/small screen banner
    echo ""
    echo -e "${ACCENT}🐿️  Squix's Live Demo!${NC}"
    echo -e "${MUTED}The SQL explorer for Terminal Squirrels${NC}"
    echo ""
else
    # Full ASCII art banner for desktop
    echo -e "${ACCENT}                                                gg${NC}"
    echo -e "${ACCENT}                                                ""${NC}"
    echo -e "${ACCENT}   n_n   __        ,g,      ,gggg,gg gg      gg  gg     ,gg,   ,gg${NC}"
    echo -e "${ACCENT}   o o  /  \      ,8'8,    dP\"  \"Y8I I8      8I  88    d8\"\"8b,dP\"${NC}"
    echo -e "${ACCENT}   |.| |  @|     ,8'  Yb  i8'    ,8I I8,    ,8I  88   dP   ,88\"${NC}"
    echo -e "${ACCENT}   \\_/\\|  _/    ,8'_   8),d8,   ,d8b,d8b,  ,d8b_,88,,dP  ,dP\"Y8,${NC}"
    echo -e "${ACCENT}  /     \\       P' \"YY8P8P\"Y8888P\"888P'\"Y88P\"'Y8P\"\"Y8\"  dP\"   \"Y8${NC}"
    echo -e "${ACCENT} /|  /  |                         I8${NC}"
    echo -e "${ACCENT}   \\__|_/                         I8'${NC}"
    echo -e "${ACCENT}  _/ _/                           I8${NC}"
    echo -e "${ACCENT}                                  I8${NC}"
    echo -e "${ACCENT}                                  I8${NC}"
    echo -e "${ACCENT}                                  I8${NC}"
    echo ""
    echo -e "${PRIMARY}                   Welcome to Squix's Live Demo!${NC}"
    echo -e "${MUTED}             =========================================${NC}"
    echo -e "${PRIMARY}              The SQL explorer for Terminal Squirrels${NC}"
    echo ""
fi
echo -e "this is a bash shell with ${ACCENT}squix${NC} and a sqlite database connection available for testing"
echo ""
echo -e "${PRIMARY}▶ Try the Demo Presentation${NC}"
echo -e "Run ${ACCENT}demo${NC} to open a live presentation showcasing squix features."
echo -e "Controls: ${ACCENT}ctrl+e${NC} run snippets • ${ACCENT}space/arrows/hjkl${NC} navigate • ${ACCENT}q/ctrl+c${NC} exit"
echo ""
echo -e "${PRIMARY}▶ Explore with Squix Commands${NC}"
echo -e "We have a demo SQLite database connected with squix with random squirrel data. Try:"
echo -e "  ${ACCENT}squix explore${NC}                To browse the database"
echo -e "  ${ACCENT}squix explain squirrels${NC}      Understand the schema"
echo -e "  ${ACCENT}squix list${NC}                   List saved queries"
echo -e "  ${ACCENT}squix add <query-name> <SQL>${NC} List saved queries"
echo -e "  ${ACCENT}squix run <query-name>${NC}       Execute a query"
echo -e "  ${ACCENT}squix help${NC}                   Get general help"
echo -e "  ${ACCENT}squix help <command>${NC}         Command-specific help"
echo ""
echo -e "${PRIMARY}▶ Recovery${NC}"
echo -e "You have write permissions to the DB. If something breaks, run ${ACCENT}reset-demo${NC} to restore the environment"
echo ""
echo -e "${PRIMARY}▶ Repo (Detailed usage and install instructions)${NC}"
echo -e "${ACCENT}http://github.com/eduardofuncao/squix${NC}"
echo ""

# Start user in tree directory (this works because .bashrc is sourced before rbash restrictions)
cd ~/tree 2>/dev/null || true

# Verify squix connection
if [[ -f ~/tree/squirrels.db ]]; then
    echo -e "${SUCCESS}✓ Squix configured with squirrels database${NC}"
fi

# Auto-recovery: restore critical files if deleted
restore_if_needed() {
    if [[ ! -f ~/tree/squirrels.db ]]; then
        cp /root/tree-templates/squirrels.db ~/tree/
        echo -e "${ERROR}Restored squirrels.db${NC}"
    fi
    if [[ ! -f ~/tree/squirrels.sql ]]; then
        cp /root/tree-templates/squirrels.sql ~/tree/
    fi
    if [[ ! -f ~/tree/squix-tutorial.md ]]; then
        cp /root/tree-templates/squix-tutorial.md ~/tree/
    fi
}

# Run recovery on shell start
if [[ $- == *i* ]]; then
    restore_if_needed
fi

# Simple prompt with squirrel emoji (no $)
export PS1='\[\e[38;5;86m\]scrat\[\e[0m\]@\[\e[38;5;205m\]squix\[\e[0m\]🐿️  '

# Alias for demo
alias demo='presenterm -x ~/tree/squix-tutorial.md'
