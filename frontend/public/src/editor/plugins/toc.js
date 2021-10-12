import {Plugin} from "prosemirror-state";

export function tocPlugin(onUpdate) {
    return new Plugin({
        state: {
            init() {},
            apply(tr) {
                if(tr.docChanged && tr.steps) {
                    for(let i = 0; i < tr.steps.length; i++) {
                        if(!tr.steps[i].slice || !tr.steps[i].slice.content || !tr.steps[i].slice.content.content) {
                            continue;
                        }

                        for(let j = 0; j < tr.steps[i].slice.content.content.length; j++) {
                            if(!tr.steps[i].slice.content.content[j].type) {
                                continue;
                            }

                            if(tr.steps[i].slice.content.content[j].type.isBlock) {
                                onUpdate({toc: true});
                            }
                        }
                    }
                }
            }
        }
    });
}
