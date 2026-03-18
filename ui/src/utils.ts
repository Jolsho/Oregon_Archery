export function random_uint32():number {
    return crypto.getRandomValues(new Uint32Array(1))[0];
}

interface HasTime {
    created_at: Date;
}

export function insertSorted<T extends HasTime>(arr: T[], item: T) {
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
