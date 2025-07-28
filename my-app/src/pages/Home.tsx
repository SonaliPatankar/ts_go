// import React, { useEffect, useState,useReducer } from 'react';
// import { Note, getNotes, createNote, updateNote, deleteNote } from '../services/api';
// import {noteReducer,initialState,State} from '../reducers/noteReducer'
// import NoteItem from '../componenets/NoteItem';

// const Home: React.FC = () => {
//   // const [notes, setNotes] = useState<Note[]>([]);
//   // const [newNote, setNewNote] = useState('');
//   // const [editId, setEditId] = useState<number | null>(null);

//   const [state, dispatch] = useReducer(noteReducer, initialState);

//   const fetchNotes = async () => {
//     const data = await getNotes();
//     dispatch({ type: 'SET_NOTES', payload: data });
//     //setNotes(data);
//   };

//   useEffect(() => {
//     fetchNotes();
//   }, []);

//   const handleSubmit = async () => {
//     if (state.editId !== null) {
//       await updateNote(state.editId, state.newNote);
//     } else {
//       await createNote(state.newNote);
//     }
//     dispatch({ type: 'RESET_INPUT' });
//     // setNewNote('');
//     // setEditId(null);
//     fetchNotes();
//   };

//   const handleEdit = (note: Note) => {
//     dispatch({type:'SET_INPUT',payload: note.content});
//     dispatch({ type: 'SET_EDIT_ID', payload: note.id });

//     // setNewNote(note.content);
//     // setEditId(note.id);
//   };

//   const handleDelete = async (id: number) => {
//     await deleteNote(id);
//     fetchNotes();
//   };

//   return (
//     <div style={{ padding: 20 }}>
//       <h1>📝 GoNotes</h1>

//       <input
//         value={state.newNote}
//        // onChange={(e) => setNewNote(e.target.value)}
//        onChange={(e) => dispatch({ type: 'SET_INPUT', payload: e.target.value })}
//         placeholder="Write a note..."
//       />
//       <button onClick={handleSubmit}>
//         {/* {editId !== null ? 'Update Note' : 'Add Note'} */}
//         {state.editId !== null ? 'Update Note' : 'Add Note'}
//       </button>

//       <ul>
//         {state.notes?.map((note) => (
//           <NoteItem
//             key={note.id}
//             note={note}
//             onEdit={handleEdit}
//             onDelete={handleDelete}
//           />
//         ))}
//       </ul>
//     </div>
//   );
// };

// export default Home;


import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { RootState,AppDispatch } from '../redux/store';
import { fetchNotes, createNote, updateNote, deleteNote, setNewNote, setEditId } from '../redux/noteSlice';
import { useEffect } from 'react';
import NoteItem from '../componenets/NoteItem';

const Home: React.FC =()=>{
  const { notes, newNote, editId, loading } = useSelector((state: RootState) => state.notes);
  const dispatch = useDispatch<AppDispatch>();
  useEffect(() => {
    dispatch(fetchNotes());
  }, [dispatch]);

  useEffect(() => {
    dispatch(fetchNotes());
  }, [dispatch]);

  const handleSubmit = () => {
    if (editId !== null) {
      dispatch(updateNote({ id: editId, content: newNote }));
    } else {
      dispatch(createNote(newNote));
    }
    dispatch({ type: 'notes/resetInput' });
  };

  const handleEdit = (note: { id: number; content: string }) => {
    dispatch(setNewNote(note.content));
    dispatch(setEditId(note.id));
  };

  const handleDelete = (id: number) => {
    dispatch(deleteNote(id));
  };
  
  return (
     <div style={{ padding: 20 }}>
      <h1>📝 GoNotes (Redux)</h1>
      <input
      
        value={newNote}
        onChange={(e) => dispatch(setNewNote(e.target.value))}
        placeholder="Write a note..."
      />
      <button onClick={handleSubmit}>
        {editId !== null ? 'Update Note' : 'Add Note'}
      </button>
      { loading ? (
        <p>Loading...</p>
      ) : (
        <ul>
          { Array.isArray(notes) ?notes.map((note) => (
            <NoteItem
              key={note.id}
              note={note}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          )): <p>No notes found.</p>}
        </ul>
      )}
    </div>
  )
}
export default Home;