import { Line } from 'react-chartjs-2';

function LineChart({ id, chartData }) {
  return (
    <Line
      datasetIdKey={id}
      data={chartData}
    />
  );
}

export default LineChart;
