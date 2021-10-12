export function scrollTo(element, padding) {
    let top = element.getBoundingClientRect().top+window.scrollY-padding;
    window.scrollTo({top, behavior: "smooth"});
}
