import { createBrowserRouter } from "react-router-dom";

const Test = () => (
  <div style={{ color: "white", padding: 20 }}>
    APP IS WORKING
  </div>
);

export const router = createBrowserRouter([
  { path: "*", element: <Test /> }
]);