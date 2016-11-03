let nextParamId = 0;
export const addParam = (name) => ({
  type: 'ADD_PARAM',
  id: nextParamId++,
  name,
});

export const removeParam = (id) => ({
  type: 'REMOVE_PARAM',
  id,
});

let nextLabelId = 0;
export const createLabel = (name) => ({
  type: 'CREATE_LABEL',
  id: nextLabelId++,
  name,
});

export const addLabel = (id) => ({
  type: 'ADD_LABEL',
  id,
});

export const removeLabel = (id) => ({
  type: 'REMOVE_LABEL',
  id,
})
