// src/redux/notesSlice.ts

import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import axios from 'axios';

export type Note = {
  id: number;
  content: string;
};

type NotesState = {
  notes: Note[];
  newNote: string;
  editId: number | null;
  loading: boolean;
  error: string | null;
};

const initialState: NotesState = {
  notes: [],
  newNote: '',
  editId: null,
  loading: false,
  error: null,
};

const BASE_URL = "https://apl4vp40xc.execute-api.us-east-1.amazonaws.com/Prod";

// Async thunks
export const fetchNotes = createAsyncThunk('notes/fetchNotes', async () => {
  const res = await axios.get<Note[]>(`${BASE_URL}/notes`, {
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  });

  return res.data;
});

export const createNote = createAsyncThunk('notes/createNote', async (content: string) => {
  const res = await axios.post<Note>(`${BASE_URL}/notes`, { content }, {
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  });
  return res.data;
});

export const updateNote = createAsyncThunk(
  'notes/updateNote',
  async ({ id, content }: { id: number; content: string }) => {
    const res = await axios.put<Note>(`${BASE_URL}/notes/${id}`, { content }, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`,
      },
    });
    return res.data;
  }
);

export const deleteNote = createAsyncThunk('notes/deleteNote', async (id: number) => {
  await axios.delete(`${BASE_URL}/notes/${id}`, {
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  });
  return id;
});

// Slice
const notesSlice = createSlice({
  name: 'notes',
  initialState,
  reducers: {
    setNewNote(state, action: PayloadAction<string>) {
      state.newNote = action.payload;
    },
    setEditId(state, action: PayloadAction<number | null>) {
      state.editId = action.payload;
    },
    resetInput(state) {
      state.newNote = '';
      state.editId = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchNotes.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchNotes.fulfilled, (state, action) => {
        state.loading = false;
        state.notes = action.payload;
      })
      .addCase(fetchNotes.rejected, (state) => {
        state.loading = false;
        state.error = 'Failed to fetch notes';
      })
      .addCase(createNote.fulfilled, (state, action) => {
        if (!Array.isArray(state.notes)) {
          state.notes = [];
        }
        state.notes.push(action.payload);
      })
      .addCase(updateNote.fulfilled, (state, action) => {
        const index = state.notes.findIndex((n) => n.id === action.payload.id);
        if (index !== -1) {
          state.notes[index] = action.payload;
        }
      })
      .addCase(deleteNote.fulfilled, (state, action) => {
        state.notes = state.notes.filter((note) => note.id !== action.payload);
      });
  },
});

export const { setNewNote, setEditId, resetInput } = notesSlice.actions;
export default notesSlice.reducer;
