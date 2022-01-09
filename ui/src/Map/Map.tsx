import React from 'react';
import { MapContainer, Marker, TileLayer } from 'react-leaflet';
import { Icon, Point} from 'leaflet';
import {Measurement, measurementRank} from '../Measurement/Measurement';
import redimg from '../img/red.png';
import blueimg from '../img/blue.png';
import orangeimg from '../img/orange.png';
import yellowimg from '../img/yellow.png';

export interface MapData {
  measurements: Measurement[]
}

type MapProps = {
  mapData: MapData,
  onMarkerClick: (m: Measurement) => void,
}

var LeafIcon = Icon.extend({
  options: {
    iconSize:     [38, 38],
    shadowSize:   [50, 64],
    iconAnchor:   [22, 94],
    shadowAnchor: [4, 62],
    popupAnchor:  [-3, -76]
  }
});

function getIcon(rank: number):Icon {
  // @ts-ignore
  let red = new LeafIcon({iconUrl: redimg});
  // @ts-ignore
  let orange = new LeafIcon({iconUrl: orangeimg});
  // @ts-ignore
  let yellow = new LeafIcon({iconUrl: yellowimg});
  // @ts-ignore
  let blue = new LeafIcon({iconUrl: blueimg});

  if (rank < 0.25) {
    return blue;
  } else if (rank < 0.5) {
    return yellow;
  } else if (rank < 0.75) {
    return orange;
  }
  return red;
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
            <Marker icon={getIcon(measurementRank(measurement))} position={[measurement.position.lat, measurement.position.long]} eventHandlers={{
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
