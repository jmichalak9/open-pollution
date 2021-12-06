export interface Measurement {
  readonly position: Position;
  readonly timestamp: Date;
  readonly temperature?: number;
  readonly levelPM10?: number;
  readonly levelPM25?: number;
  readonly levelSO2?: number;
  readonly levelO3?: number;

}

export interface Position {
  readonly lat: number;
  readonly long: number;
}