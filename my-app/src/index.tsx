import React from 'react';
import ReactDOM from 'react-dom/client';
import reportWebVitals from './reportWebVitals';
import { Provider } from 'react-redux';
import { store } from './redux/store';
import { BrowserRouter } from 'react-router-dom';
// import { Security } from '@okta/okta-react';
import './index.css';
import App from './App';
import { oktaConfig } from './componenets/oktaConfig';
import { OktaAuth } from '@okta/okta-auth-js';

const oktaAuth = new OktaAuth(oktaConfig);

const restoreOriginalUri = async (_oktaAuth: any, originalUri: string) => {
  window.location.replace(originalUri || '/');
};

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  // <Security oktaAuth={oktaAuth} restoreOriginalUri={restoreOriginalUri}>
    <Provider store={store}>
     <BrowserRouter>
      <App />
    </BrowserRouter>
  </Provider>
  // {/* </Security> */}
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
