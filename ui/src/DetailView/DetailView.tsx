import React from 'react';
import {Measurement} from "../Measurement/Measurement";

export interface MapData {
  measurement: Measurement
}

type DetailViewProps = {
  mapData: MapData,
}

const DetailView: React.FC<DetailViewProps> = ({ mapData }: DetailViewProps) => {
  const lat = mapData.measurement.position.lat;
  const long = mapData.measurement.position.long;

  return (
    <div className="DetailView">
      <h1>Szczegóły pomiaru</h1>
      <div className="LeftDetails">
        <div>Szerokość geograficzna {lat}°{lat > 0 ? <text>N</text> : <text>S</text>}</div>
        <div>Długość geograficzna {long}°{long > 0 ? <text>E</text> : <text>W</text>}</div>
        { mapData.measurement.temperature ? <div>Temperatura {mapData.measurement.temperature}°C</div> : <div></div>}
      </div>
      <div className="RightDetails">
        { mapData.measurement.levelPM10 ? <div>Zanieczyszczenie PM 10 {mapData.measurement.levelPM10} μg/m<sup>3</sup></div> : <div></div>}
        { mapData.measurement.levelPM25 ? <div>Zanieczyszczenie PM 2.5 {mapData.measurement.levelPM25} μg/m<sup>3</sup></div> : <div></div>}
        { mapData.measurement.levelSO2 ? <div>Zanieczyszczenie SO<sub>2</sub> {mapData.measurement.levelSO2} μg/m<sup>3</sup></div> : <div></div>}
        { mapData.measurement.levelO3 ? <div>Zanieczyszczenie O<sub>3</sub> {mapData.measurement.levelO3} μg/m<sup>3</sup></div> : <div></div>}
      </div>
    </div>
  );
}

export default DetailView;
