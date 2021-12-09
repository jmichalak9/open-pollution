import React from 'react';
import ReactDOM from 'react-dom';
import { positions, types, Provider as AlertProvider } from 'react-alert';
import AlertTemplate from './Alert/Alert';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

// optional configuration
const options = {
    position: positions.BOTTOM_CENTER,
    timeout: 5000,
};

ReactDOM.render(
  <React.StrictMode>
      <AlertProvider template={AlertTemplate} {...options}>
          <App />
      </AlertProvider>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
