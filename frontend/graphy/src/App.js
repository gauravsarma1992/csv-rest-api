// App.js
import { useState, useEffect } from "react";
import axios from 'axios';
import Dashboard from "./components/Dashboard";

export default function App() {
  return (
    <div className="App">
      <Dashboard />
    </div>
  );
}
