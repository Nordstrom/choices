/**
 * shuffle shuffles an array using fisher-yates algorithm
 * @param {Array} arr - the array to shuffle.
 */
export function shuffle(arr) {
  for (let i = arr.length-1; i > 0; i--) {
    const n = Math.floor(Math.random()*i);
    [arr[i], arr[n]] = [arr[n], arr[i]];
  }
  return arr
}

export function sample(arr1, num) {
  const free = arr1.map((seg, i) => {
    // if the segment is claimed return -1
    if (seg === 1) {
      return -1;
    }
    // otherwise return the index
    return i;
  }).filter(i => i >= 0);
  shuffle(free);
  let out = new Array(128).fill(0);
  free.slice(0, num).forEach(i => out[i] = 1);
  return out;
}