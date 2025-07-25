import React from 'react';
import { Note } from '../services/api';

type Props = {
  note: Note;
  onEdit: (note: Note) => void;
  onDelete: (id: number) => void;
};

const NoteItem: React.FC<Props> = ({ note, onEdit, onDelete }) => {
  return (
    <li>
      {note.content}

      <button onClick={() => onEdit(note)}>Edit</button>
      <button onClick={() => onDelete(note.id)}>Delete</button>
    </li>
  );
};

export default NoteItem;
