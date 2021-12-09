import React, {useEffect, useState} from 'react';
import { useAlert } from 'react-alert';

import LeftPanel from './LeftPanel/LeftPanel';
import Map, {MapData} from './Map/Map';

import './App.css';
import DetailView from "./DetailView/DetailView";
import {Measurement} from "./Measurement/Measurement";
import {getMeasurements} from "./APIClient/APIClient";


const mockDetails = {
  measurement: {} as Measurement,
};

function App() {
  const alert = useAlert();
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
      }, () => {
        alert.show('Nie udało się pobrać danych!')
      });
    };
    return x;
  }
  useEffect(() => {
    getMeasurementsClosure();
    setInterval(getMeasurementsClosure(), 3000);
  },[]);

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
