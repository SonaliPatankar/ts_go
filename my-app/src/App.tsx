import React from 'react';
import logo from './logo.svg';
import './App.css';
import MyButton from './componenets/button';
import { CounterProvider } from './context/counter';
import Home from './pages/Home';
import ApiCall from './componenets/apiCall';
import { Provider } from 'react-redux';
import { store } from './redux/store';
import Login from './pages/Login';
import { Route, Routes } from 'react-router';
function App() {
  return (
    // <CounterProvider>
      <div className="App">
        {/* <MyButton text="Click Me" />
        <MyButton onclick={() => alert("Another Button Clicked")} text="Another Button" /> */}
        {/* <MyButton text="Increment Count" /> */}
<Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/" element={<Home />} />
    </Routes>
      </div>
    
  );
}

export default App;
