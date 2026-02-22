
export function random_uint32() {
    return crypto.getRandomValues(new Uint32Array(1))[0];
}

/** 
 * @param {string} kind 
 * @param {HTMLElement | null} parent 
 * @param {...string} classes 
 */
export function create_elem(kind, parent, ...classes) {
    let elem = document.createElement(kind);
    if (classes.length) elem.classList.add(...classes);
    if (parent) parent.appendChild(elem);
    return elem;
}
