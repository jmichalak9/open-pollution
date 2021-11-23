import React, {useEffect, useState} from 'react';

import LeftPanel from './LeftPanel/LeftPanel';
import Map, {MapData} from './Map/Map';

import logo from './logo.svg';
import './App.css';
import DetailView from "./DetailView/DetailView";
import {Measurement} from "./Measurement/Measurement";
import {getMeasurements} from "./APIClient/APIClient";

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
  const [measurements, setMeasurements] = useState<MapData>( {
    measurements: [] as Measurement[]
  });

  function updateDetailsView(measurement: Measurement) {
    const details = {
      measurement: measurement,
    };
    setDetails(details);
    console.log(measurement);
  }
  function getMeasurementsClosure(): Function {
    let x = ()=> {
      getMeasurements((m: Measurement[]) => {
        console.log("APP.tsx: ", m);
        const measurements = {
          measurements: m,
        }
        setMeasurements(measurements);
      });
    };
    return x;
  }
  useEffect(() => {
    getMeasurementsClosure();
    setInterval(getMeasurementsClosure(), 3000);
  },[]);

  async function loadData() {
    try {
      const res = await fetch('https://example.com');
      const blocks = await res;//.json();
      console.log(blocks)
    } catch (e) {
      console.log(e);
    }
  }
  return (
    <div className="App">
      <LeftPanel/>
      <div className="Container">
        <div className="CenterElements">
          <Map mapData = {measurements} onMarkerClick={updateDetailsView}/>
          <div className="Details">
            <h1>Szczegóły pomiaru</h1>
            <hr />
            <DetailView mapData={details}/>
          </div>
        </div>
      </div>
    </div>

  );
}


export default App;
