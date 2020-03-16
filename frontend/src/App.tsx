import React from "react";

import "./App.scss";

import Something from "./something/Something";
import { GetSomething } from "./something/SomethingAPI";

const App = () => {
  return (
    <div className="App">
      <Something getSomething={GetSomething} />
    </div>
  );
};

export default App;
