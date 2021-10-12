// getCmdType returns the type of given command string.
// This can either be "search", "command" or "navigate" depending on the first character of the string.
// Or null if it's empty.
export function getCmdType(cmd) {
    if(cmd === "") {
        return null;
    }

    if(cmd.startsWith("/")) {
        return "command";
    }

    if(cmd.startsWith(".")) {
        return "navigation";
    }

    return "search";
}
