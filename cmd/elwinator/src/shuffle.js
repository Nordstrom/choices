/**
 * shuffle shuffles an array using fisher-yates algorithm
 * @param {Array} arr - the array to shuffle.
 */
export default function shuffle(arr) {
  for (let i = arr.length-1; i > 0; i--) {
    const n = Math.floor(Math.random()*i);
    [arr[i], arr[n]] = [arr[n], arr[i]];
  }
  return arr
}