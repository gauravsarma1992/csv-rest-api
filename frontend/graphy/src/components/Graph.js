// App.js
import { useState, useEffect } from "react";
import { Data } from "../utils/Data";
import LineChart from "./LineChart";
import { getUrl } from "./common";
import axios from 'axios';

import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  BarElement,
} from 'chart.js';
import { Chart } from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  BarElement,
  ArcElement
)

function generateGraphData(payload, matched_elems, key) {
  let labels = matched_elems.map((data) => data.time);
  let datasets = [{
    label: key.toUpperCase() + " for API - " + payload.filters.requestname,
    data: matched_elems.map((data) => data[key]),
    borderColor: 'rgb(53, 162, 235)',
    backgroundColor: 'rgba(53, 162, 235, 0.5)',
  }];
  return {
    labels: labels,
    datasets: datasets,
  };
}

export default function Graph({ files, filters}) {
  const sp_file_name = files[0].split("-");
  const experimentName = "Provider - " + sp_file_name[sp_file_name.length-2] + " | Action - " + sp_file_name[sp_file_name.length-1].split(".")[0]; 
  const [latencyData, setLatencyData] = useState({
    labels: [],
    datasets: {},
  });
  const [statusData, setStatusData] = useState({
    labels: [],
    datasets: {},
  });

  useEffect(() => {
    var payload = {
      files: files,
      filters: filters 
    }
    axios.post(getUrl("csv"), payload).then((response) => {
      const matched_elems = response.data.matched_elems;
      setLatencyData(generateGraphData(payload, matched_elems, "latency"));
      setStatusData(generateGraphData(payload, matched_elems, "statuscode"));
    });
  }, []);

  return (
    <Container md="fluid">
      <h3>Latency and status graph</h3> 
      <h4>{experimentName} | {filters.requestname} </h4>
      {latencyData.labels.length > 0 ? 
      <Row>
        <Col md={12}>
          <LineChart id="temp" chartData={latencyData} />
        </Col>
      </Row>: null }
      {statusData.labels.length > 0 ? 
      <Row>
        <Col md={12}>
          <LineChart id="temp" chartData={statusData} />
        </Col>
      </Row>: null }
    </Container>
  );
}
