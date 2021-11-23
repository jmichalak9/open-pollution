import React from 'react';
import { MapContainer, Marker, TileLayer } from 'react-leaflet';
import { Measurement} from '../Measurement/Measurement';

export interface MapData {
  measurements: Measurement[]
}

type MapProps = {
  mapData: MapData,
  onMarkerClick: (m: Measurement) => void,
}

const Map: React.FC<MapProps> = ({ mapData, onMarkerClick }: MapProps) => {
  console.log(mapData);
    return (
      <div className="Map">
        <MapContainer center={[52.2297, 21.0122]} zoom={13}>
          <TileLayer
            attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          {mapData.measurements.map((measurement) =>
            <Marker position={[measurement.position.lat, measurement.position.long]} eventHandlers={{
                click: () => {
                  onMarkerClick(measurement);
                },
            }}>
            </Marker>
          )}
        </MapContainer>
      </div>
    );
}

export default Map;
