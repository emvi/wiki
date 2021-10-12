const scrollDistance = 48; // distance to jump on scroll
const scrollThreshold = 80*2;
const scrollThresholdArea = 48;
const scrollTopAreaPadding = 8;

// scroll scrolls the whole page.
export function scroll(dir) {
    window.scroll(window.scrollX, window.scrollY+scrollDistance*dir);
}

// scrollArea scrolls given area up or down, depending on direction (> 0 down, < 0 up).
export function scrollArea(area, dir) {
    area.scrollTop = area.scrollTop+scrollDistance*dir;
}

// scrollIntoView scrolls given element into the visible area of the window with some padding to top and bottom.
// The padding is the delta between top/bottom distance and the scrollThreshold.
export function scrollIntoView(element) {
    let height = window.innerHeight;
    let box = element.getBoundingClientRect();
    let distanceTop = box.top;
    let distanceBottom = height-box.bottom;

    if(distanceTop < scrollThreshold) {
        window.scroll(window.scrollX, window.scrollY-(scrollThreshold-distanceTop));
    }
    else if(distanceBottom < scrollThreshold) {
        window.scroll(window.scrollX, window.scrollY+(scrollThreshold-distanceBottom));
    }
}

// scrollIntoViewArea does the same as scrollIntoView but for given area.
export function scrollIntoViewArea(element, area) {
    let elementBox = element.getBoundingClientRect();
    let areaBox = area.getBoundingClientRect();
    let height = areaBox.height;
    let distanceTop = elementBox.top-areaBox.top;
    let distanceBottom = height-(elementBox.bottom-areaBox.top);

    if(distanceTop < scrollThresholdArea) {
        area.scrollTop = area.scrollTop-(scrollThresholdArea-distanceTop);
    }
    else if(distanceBottom < scrollThresholdArea) {
        area.scrollTop = area.scrollTop+(scrollThresholdArea-distanceBottom);
    }
}

// scrollToTop scrolls given element to the top most position minus scrollThreshold.
// This can be used if an element takes more space to the bottom in certain scenarios (like the preview).
export function scrollToTop(element) {
    let box = element.getBoundingClientRect();
    window.scroll(window.scrollX, window.scrollY+box.top-scrollThreshold);
}

// scrollToTopArea does the same as scrollToTop but for given area.
export function scrollToTopArea(element, area) {
    let elementBox = element.getBoundingClientRect();
    let areaBox = area.getBoundingClientRect();
    area.scrollTop = area.scrollTop+(elementBox.top-areaBox.top)-scrollTopAreaPadding;
}

export function scrollTo(element, padding) {
    let top = element.getBoundingClientRect().top+window.scrollY-padding;
    window.scrollTo({top, behavior: "smooth"});
}
