export function copyToClipboard(value) {
    navigator.permissions.query({name: "clipboard-write"})
    .then(result => {
        if(result.state === "granted" || result.state === "prompt") {
            navigator.clipboard.writeText(value);
        }
    })
    .catch(() => {
        // lets do it anyways!
        navigator.clipboard.writeText(value);
    });
}
