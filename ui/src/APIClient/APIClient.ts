import axios from "axios";

const backendURL = "http://localhost:8000";
const measurementsPath = "/measurements";

type measurementsAPIResponse = {
  // TODO: fill this.
};

export async function getMeasurements(): Promise<measurementsAPIResponse> {
  const response = await axios.get(backendURL + measurementsPath);

  return Promise.resolve(response.data["texts"]).catch(err =>
    alert(err)
  );
}