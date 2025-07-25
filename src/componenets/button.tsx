import { title } from "process";
import React, { use } from "react"
import { useState } from "react";
import  { useCounter } from "../context/counter";

interface MyButtonProps {
  text: string;
  onclick?: () => void;
}

// const MyButton: React.FC<MyButtonProps> = (props) => {
    
//   return <button onClick={props.onclick}>{props.text}</button>;
// }
// const MyButton: React.FC<MyButtonProps> = (props) => {
    
//   const { text, onclick } = props;
//   return (
//     <button onClick={onclick}>
//       {text}
//     </button>
//   );    
// }

interface Book {
  title: string;
  author: string;
  price: number;
}

// const MyButton: React.FC<MyButtonProps> = (props) => {
    
//   const { text, onclick } = props;
//   const [count, setCount] = useState<Book>({
//     title: "Default Title",
//     author: "Default Author",
//     price: 10
//   });
//   return (
//     <>
//       <h3>{count.title} by {count.author} - ${count.price}</h3>
//       <button onClick={()=> setCount({title:"updated title", author:"updated author", price: 20})}>
//         {text}
//       </button>
//     </>
//   );    
// }

// const MyButton: React.FC<MyButtonProps> = (props) => {
    
//   const { text, onclick } = props;
//   const [count, setCount] = useState<Book>({
//     title: "Default Title",
//     author: "Default Author",
//     price: 10
//   });
//   const [value, setValue] = useState<string | undefined>();

//   const handleChange=(e:React.ChangeEvent<HTMLInputElement>)=>{
//     setValue(e.target.value);
//     console.log("Input Value:", e.target.value);
//   }

//   const handleSubmit=(e:React.FormEvent<HTMLFormElement>)=>{
//     e.preventDefault();
//     console.log("Form Submitted with Value:", value);
//   }
//   return (
    
//     <div>
//         <form onSubmit={handleSubmit}>
//         <input 
//           type="text" 
//           onChange={handleChange}
//         />
//         <button type="submit">Submit</button>
//         <h1>{value} </h1>
//         </form>
//     </div>
//   );    
// }

const MyButton: React.FC<MyButtonProps> = (props) => {

    const context = useCounter();
    
  return (
    <>
      <h1 onClick={() => context?.setCount(context.value + 1)}>
        Count: {context?.value}
      </h1>
    </>
  );    
}


export default MyButton;