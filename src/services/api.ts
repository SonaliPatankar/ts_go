import axios from 'axios';
import { create } from 'domain';

const API_URL = 'http://localhost:8080/notes';

export type Note = {
  id: number;
  content: string;
};


export const getNotes = async (): Promise<Note[]> => {
  try {
    const res = await axios.get<Note[]>(API_URL);
    return res.data;
  } catch (error: any) {
    console.error('Failed to fetch notes:', error);
    throw new Error(error?.response?.data?.message || 'Failed to fetch notes');
  }
};


export const createNote = async (content: string): Promise<Note> => {
  try {
    const res = await axios.post<Note>(API_URL, { content });
    console.log('Created Note:', res.data);
    return res.data;
  } catch (error: any) {
    console.error('Failed to create note:', error);
    throw new Error(error?.response?.data?.message || 'Failed to create note');
  }
};


export const updateNote = async (id: number, content: string): Promise<Note> => {
  try {
    const res = await axios.put<Note>(`${API_URL}/${id}`, { content });
    return res.data;
  } catch (error: any) {
    console.error('Failed to update note:', error);
    throw new Error(error?.response?.data?.message || 'Failed to update note');
  }
};


export const deleteNote = async (id: number): Promise<void> => {
  try {
    await axios.delete(`${API_URL}/${id}`);
  } catch (error: any) {
    console.error('Failed to delete note:', error);
    throw new Error(error?.response?.data?.message || 'Failed to delete note');
  }
};


