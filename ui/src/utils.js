
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

export function insertSorted(arr, item) {
  const t = new Date(item.created_at).getTime();
  let lo = 0, hi = arr.length;

  while (lo < hi) {
    const mid = (lo + hi) >> 1;
    const mt = new Date(arr[mid].created_at).getTime();

    if (mt <= t) lo = mid + 1;
    else hi = mid;
  }

  arr.splice(lo, 0, item);
  return arr;
}
