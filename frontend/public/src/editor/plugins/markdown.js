import {InputRule, inputRules, wrappingInputRule, textblockTypeInputRule} from "prosemirror-inputrules";

function headingRule(nodeType, maxLevel) {
    return textblockTypeInputRule(new RegExp("^(#{1," + maxLevel + "})\\s$"), nodeType, match => {
        return {level: match[1].length+1};
    });
}

function blockQuoteRule(nodeType) {
    return wrappingInputRule(/^\s*>\s$/, nodeType);
}

function orderedListRule(nodeType) {
    return wrappingInputRule(/^(\d+)\.\s$/, nodeType, match => {
        return {order: match[1]};
    }, (match, node) => {
        return node.childCount + node.attrs.order == +match[1];
    });
}

function bulletListRule(nodeType) {
    return wrappingInputRule(/^\s*([-+*])\s$/, nodeType);
}

function uncheckedListRule(nodeType) {
    return wrappingInputRule(/^(\[[ ]])\s$/, nodeType);
}

function checkedListRule(nodeType) {
    return wrappingInputRule(/^(\[[x\-Oo]])\s$/, nodeType, () => {
        return {checked: true};
    });
}

function codeBlockRule(nodeType) {
    return textblockTypeInputRule(/^```$/, nodeType);
}

function horizontalLineRule(regexp, nodeType) {
    return new InputRule(regexp , (state, match, start, end) => {
        let tr = state.tr;
        tr.delete(start, end);
        tr.replaceSelectionWith(nodeType.create());
        return tr;
    });
}

function boldItalicInputRule(regexp, boldMark, italicMark) {
    return new InputRule(regexp , (state, match, start, end) => {
        if(!match[2]) {
            return;
        }

        let markType = boldMark;

        if(match[1] === "*") {
            markType = italicMark;
        }

        let tr = state.tr;
        let textStart = start + match[0].indexOf(match[2]);
        let textEnd = textStart + match[2].length
        if (textEnd < end) tr.delete(textEnd, end);
        if (textStart > start) tr.delete(start, textStart);
        end = start + match[2].length;
        tr.addMark(start, end, markType.create());
        tr.removeStoredMark(markType);
        return tr;
    });
}

function markInputRule(regexp, markType) {
    return new InputRule(regexp, (state, match, start, end) => {
        let tr = state.tr;

        if(match[1]) {
            let textStart = start + match[0].indexOf(match[1]);
            let textEnd = textStart + match[1].length
            if (textEnd < end) tr.delete(textEnd, end);
            if (textStart > start) tr.delete(start, textStart);
            end = start + match[1].length;
        }

        tr.addMark(start, end, markType.create());
        tr.removeStoredMark(markType);
        return tr;
    });
}

export function markdownPlugin(schema) {
    return inputRules({rules: [
        headingRule(schema.nodes.headline, 3),
        blockQuoteRule(schema.nodes.blockquote),
        orderedListRule(schema.nodes.ordered_list),
        bulletListRule(schema.nodes.bullet_list),
        uncheckedListRule(schema.nodes.check_list_item),
        checkedListRule(schema.nodes.check_list_item),
        codeBlockRule(schema.nodes.code_block),
        boldItalicInputRule(/(\*{1,2})(.*?)\1/s, schema.marks.bold, schema.marks.italic),
        markInputRule(/(?:`)([^`]+)(?:`)$/s, schema.marks.code),
        markInputRule(/(?:~~)([^~]+)(?:~~)$/s, schema.marks.strikethrough),
        horizontalLineRule(/^([-_*]){3,}[\t]*$/s, schema.nodes.horizontal_rule)
    ]});
}
