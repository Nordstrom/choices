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
