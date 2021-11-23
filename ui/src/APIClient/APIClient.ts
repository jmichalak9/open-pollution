import axios from "axios";
import {Measurement} from "../Measurement/Measurement";
// TODO: read this from env
const backendURL = "http://localhost:8000";
const measurementsPath = "/measurements";

interface measurementsAPIResponse {
  data: measurementFromAPI[];
};
interface measurementFromAPI {
  levelO3: number;
}

export async function getMeasurements(callback: Function) {
  // @ts-ignore
  const response = await axios.request<measurementsAPIResponse>({
    url: backendURL + measurementsPath
  }).catch(err =>
    console.log(err)
  ).then( resp => {
      // @ts-ignore
    const { data } = resp;
    let measurements: Measurement[] = [];
    for (let i = 0; i < data.length; i++) {
      const m: Measurement = {
        temperature: data[i].temperature,
        levelPM10: data[i].levelPM10,
        levelPM25:data[i].levelPM25,
        levelSO2: data[i].levelS02,
        levelO3: data[i].levelO3,
        position: {
          lat: data[i].position.lat,
          long: data[i].position.long,
        }
      };
      measurements.push(m)
    }
    callback(measurements);
    }
  );
  console.log("GETTING DATA")
}