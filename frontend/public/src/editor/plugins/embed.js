import axios from "axios";
import {Plugin} from "prosemirror-state";

let rules = [
    // YouTube
    {
        match(clipboard) {
            let regex = new RegExp(/^https?:\/\/(www\.)?youtube\.[a-zA-Z0-9()]{1,6}\b\/watch\?v=[a-zA-Z0-9-]+$/gi);
            return clipboard.match(regex);
        },
        apply(schema, view, clipboard) {
            let regex = new RegExp(/watch\?v=([a-zA-Z0-9-]+)/gi);
            let matches = regex.exec(clipboard);

            if(matches.length < 2) {
                return;
            }

            let src = `https://www.youtube.com/embed/${matches[matches.length-1]}`;
            let tr = view.state.tr;
            tr.replaceSelectionWith(schema.nodes.youtube.create({src}));
            view.dispatch(tr);
        }
    },
    // Vimeo
    {
        match(clipboard) {
            let regex = new RegExp(/^https?:\/\/(www\.)?vimeo\.[a-zA-Z0-9()]{1,6}\b\/[a-zA-Z0-9]+$/gi);
            return clipboard.match(regex);
        },
        apply(schema, view, clipboard) {
            let regex = new RegExp(/^https?:\/\/(www\.)?vimeo\.[a-zA-Z0-9()]{1,6}\b\/([a-zA-Z0-9]+)$/gi);
            let matches = regex.exec(clipboard);

            if(matches.length < 2) {
                return;
            }

            let src = `https://player.vimeo.com/video/${matches[matches.length-1]}`;
            let tr = view.state.tr;
            tr.replaceSelectionWith(schema.nodes.vimeo.create({src}));
            view.dispatch(tr);
        }
    },
    // Spotify (URL)
    {
        match(clipboard) {
            let regex = new RegExp(/^https?:\/\/(www\.)?open\.spotify\.[a-zA-Z0-9()]{1,6}\b\/(playlist|track|artist|album){1}\/[a-zA-Z0-9?=-]+$/gi);
            return clipboard.match(regex);
        },
        apply(schema, view, clipboard) {
            let regex = new RegExp(/^https?:\/\/(www\.)?open\.spotify\.[a-zA-Z0-9()]{1,6}\b\/(playlist|track|artist|album){1}\/([a-zA-Z0-9?=-]+)$/gi);
            let matches = regex.exec(clipboard);

            if(matches.length < 3) {
                return;
            }

            let type = matches[matches.length-2];
            let id = matches[matches.length-1];
            let src = `https://open.spotify.com/embed/${type}/${id}`;
            let tr = view.state.tr;
            tr.replaceSelectionWith(schema.nodes.spotify.create({src}));
            view.dispatch(tr);
        }
    },
    // Spotify (Spotify URI)
    {
        match(clipboard) {
            let regex = new RegExp(/^spotify:(playlist|track|artist|album){1}:[a-zA-Z0-9]+$/gi);
            return clipboard.match(regex);
        },
        apply(schema, view, clipboard) {
            let regex = new RegExp(/^spotify:(playlist|track|artist|album){1}:([a-zA-Z0-9]+)$/gi);
            let matches = regex.exec(clipboard);

            if(matches.length < 3) {
                return;
            }

            let type = matches[matches.length-2];
            let id = matches[matches.length-1];
            let src = `https://open.spotify.com/embed/${type}/${id}`;
            let tr = view.state.tr;
            tr.replaceSelectionWith(schema.nodes.spotify.create({src}));
            view.dispatch(tr);
        }
    },
    // picture
    {
        match(clipboard) {
            let regex = new RegExp(/^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)\.(svg|jpg|jpeg|gif|png){1}$/gi);
            return clipboard.match(regex);
        },
        apply(schema, view, clipboard) {
            let tr = view.state.tr;
            tr.replaceSelectionWith(schema.nodes.image.create({src: clipboard},
                schema.nodes.paragraph.create()));
            view.dispatch(tr);
        }
    },
    // PDF (external)
    {
        match(clipboard) {
            let regex = new RegExp(/^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)\/?[a-zA-Z0-9\-\_]\.pdf$/gi);
            return clipboard.match(regex);
        },
        apply(schema, view, clipboard) {
            let tr = view.state.tr;
            tr.replaceSelectionWith(schema.nodes.pdf.create({src: clipboard}));
            view.dispatch(tr);
        }
    },
    // link previews
    {
        match(clipboard) {
            let regex = new RegExp(/^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$/gi);
            return clipboard.match(regex);
        },
        apply(schema, view, clipboard) {
            axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/urlmeta`, {params: {url: clipboard}})
            .then(data => {
                let meta = data.data;
                let tr = view.state.tr;
                tr.replaceSelectionWith(schema.nodes.link_preview.create({href: clipboard, title: meta.title, description: meta.description, image: meta.image}));
                view.dispatch(tr);
            })
            .catch(e => {
                console.error(e);
                let tr = view.state.tr;
                let text = schema.text(clipboard, [schema.marks.link.create({href: clipboard})]);
                tr.replaceSelectionWith(schema.nodes.paragraph.create(null, text));
                view.dispatch(tr);
            });
        }
    }
];

export function embedPlugin(schema) {
    return new Plugin({
        props: {
            handlePaste(view, event) {
                let clipboard = event.clipboardData.getData("text").trim();

                if(!clipboard) {
                    return false;
                }

                let selection = view.state.selection;

                if(selection.from === selection.to && selection.$from && selection.$from.parent.type.name === "paragraph" && selection.$from.parent.content.size === 0) {
                    for(let i = 0; i < rules.length; i++) {
                        if(rules[i].match(clipboard)) {
                            rules[i].apply(schema, view, clipboard);
                            return true;
                        }
                    }
                }

                return false;
            }
        }
    });
}
