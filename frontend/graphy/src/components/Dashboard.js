import { useState, useEffect } from "react";
import axios from 'axios';
import Graph from "./Graph";
import { getUrl, getRequestname } from "./common";

export default function Dashboard() {
  const [ payloads, setPayloads ] = useState([]);
  useEffect(() => {
    const folderUrl = getUrl("folder");
    const folder_name = "/tmp/chaos-stats";
    axios.get(folderUrl).then((response) => {
      const folder_files = response.data.folder_files;
      let payload_data = [];
      folder_files.map(function(file_name) {
        payload_data.push({
          files: [folder_name + "/" + file_name],
          filters: {
            requestname: getRequestname(),
          }
        });
      }); 
      setPayloads(payload_data);
    });
  }, []);
  return (
    <div className="dashboard">
      {payloads.map(function(payload, i){
        return <Graph files={payload.files} filters={payload.filters} />
      })}
    </div>
  );
}
