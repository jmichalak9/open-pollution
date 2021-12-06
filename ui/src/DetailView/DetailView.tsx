import React from 'react';
import {Measurement} from "../Measurement/Measurement";

export interface MapData {
  measurement: Measurement
}

type DetailViewProps = {
  mapData: MapData,
}

const zeroPad = (num: number, places: number) => String(num).padStart(places, '0');

const DetailView: React.FC<DetailViewProps> = ({ mapData }: DetailViewProps) => {
  if (Object.keys(mapData.measurement).length === 0) {
    return(<div></div>)
  }
  const lat = +mapData.measurement.position.lat.toFixed(4);
  const long = +mapData.measurement.position.long.toFixed(4);
  const t = mapData.measurement.timestamp;
  return (
  <div className="DetailView">
    <div className="DetailsFlexPanel">

      <div>Szerokość geograficzna</div>
      <div>{lat}°{lat > 0 ? <text>N</text> : <text>S</text>}</div>

      <div>Długość geograficzna</div>
      <div>{long}°{long > 0 ? <text>E</text> : <text>W</text>}</div>

      <div>Czas pomiaru</div>
      <div>{zeroPad(t.getHours(), 2)}:{zeroPad(t.getMinutes(), 2)}:{zeroPad(t.getSeconds(), 2)}
        <text>     </text>{zeroPad(t.getDate(), 2)}-{zeroPad(t.getMonth()+1, 2)}-{zeroPad(t.getFullYear(), 4)}</div>

      { mapData.measurement.temperature ? <div className="DetailsFlexPanel"><div>Temperatura</div> <div>{mapData.measurement.temperature}°C</div></div> : <div></div>}
    </div>
    <div>
      { mapData.measurement.levelPM10 ? <div className="DetailsFlexPanel"><div>Zanieczyszczenie PM 10</div> <div>{mapData.measurement.levelPM10} μg/m<sup>3</sup></div></div> : <div></div>}
      { mapData.measurement.levelPM25 ? <div className="DetailsFlexPanel"><div>Zanieczyszczenie PM 2.5</div> <div>{mapData.measurement.levelPM25} μg/m<sup>3</sup></div></div> : <div></div> }
      { mapData.measurement.levelSO2 ? <div className="DetailsFlexPanel"><div>Zanieczyszczenie SO<sub>2</sub></div> <div>{mapData.measurement.levelSO2} μg/m<sup>3</sup></div></div> : <div></div>}
      { mapData.measurement.levelO3 ? <div className="DetailsFlexPanel"><div>Zanieczyszczenie O<sub>3</sub></div> <div>{mapData.measurement.levelO3} μg/m<sup>3</sup></div></div>: <div></div>}
    </div>
  </div>

);
}

export default DetailView;
