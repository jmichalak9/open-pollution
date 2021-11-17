import React, {useState} from 'react';

import LeftPanel from './LeftPanel/LeftPanel';
import Map from './Map/Map';

import logo from './logo.svg';
import './App.css';
import DetailView from "./DetailView/DetailView";
import {Measurement} from "./Measurement/Measurement";

const mockMeasurement1: Measurement = {
  temperature: 42,
  levelPM10: 33.3,
  position: {
    lat: 52.2,
    long: 21.0,
  }
};

const mockMeasurement2: Measurement = {
  temperature: 42,
  levelPM10: 33.3,
  levelPM25:30,
  levelSO2: 22,
  position: {
    lat: 52.22,
    long: 21.01,
  }
};

const mockMeasurements = {
  measurements: [mockMeasurement1, mockMeasurement2],
};

const mockDetails = {
  measurement: mockMeasurement1,
};

function App() {
  // @ts-ignore
  const [details, setDetails] = useState(mockDetails);

  function updateDetailsView(measurement: Measurement) {
    const details = {
      measurement: measurement,
    };
    setDetails(details);
    console.log(measurement);
  }
  return (
    <div className="App">
      <LeftPanel/>
      <Map mapData = {mockMeasurements} onMarkerClick={updateDetailsView}/>
      <DetailView mapData={details}/>
    </div>
  );
}


export default App;
