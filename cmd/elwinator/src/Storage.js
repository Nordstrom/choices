export const loadState = () => {
  try {
    const serializedState = localStorage.getItem('state');
    if (serializedState === null) {
      return undefined;
    }
    const s = JSON.parse(serializedState);
    s.entities.experiments = s.entities.experiments.map(e => {
      e.segments.length = 16;
      return {
        ...e,
        segments: new Uint8Array(e.segments),
      }
    });
    return s;
  } catch (err) {
    return undefined;
  }
};

export const saveState = (state) => {
  try {
    const serializedState = JSON.stringify(state);
    localStorage.setItem('state', serializedState);
  } catch (err) {
    // ignore err
  }
}