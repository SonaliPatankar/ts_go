export type Note = {
    id: number;
    content: string;
}

export type Action=
    | { type: 'SET_NOTES'; payload: Note[] }
    | { type: 'SET_INPUT'; payload: string }
    | { type: 'SET_EDIT_ID'; payload: number | null }
    | { type: 'RESET_INPUT' };

export type State={
    notes: Note[];
    newNote: string;
    editId: number | null;
}
export const initialState: State = {
  notes: [],
  newNote: '',
  editId: null,
};
export const noteReducer =(state:State, action:Action) => {
    switch (action.type) {
        case 'SET_NOTES':
            return { ...state, notes: action.payload };
        case 'SET_INPUT':
            return { ...state, newNote: action.payload };
        case 'SET_EDIT_ID':
            return { ...state, editId: action.payload };
        case 'RESET_INPUT':
            return { ...state, newNote: '', editId: null };
        default:
            return state;
    }
};
