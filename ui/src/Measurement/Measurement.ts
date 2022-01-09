export interface Measurement {
  readonly position: Position;
  readonly timestamp: Date;
  readonly temperature?: number;
  readonly levelPM10?: number;
  readonly levelPM25?: number;
  readonly levelSO2?: number;
  readonly levelO3?: number;

}

// measurementRank returns a value close to 0 if pollution level is good and close to 1 if bad.
export function measurementRank(m: Measurement): number {
  let sum = 0;
  let count = 0;
  if (m.levelPM10 !== undefined) {
    count++;
    sum += m.levelPM10;
  }
  if (m.levelPM25 !== undefined) {
    count++;
    sum += m.levelPM25;
  }
  if (m.levelSO2 !== undefined) {
    count++;
    sum += m.levelSO2;
  }
  if (m.levelO3 !== undefined) {
    count++;
    sum += m.levelO3;
  }
  if (count == 0) {
    return 0;
  }
  return sum / count / 100;
}

export interface Position {
  readonly lat: number;
  readonly long: number;
}