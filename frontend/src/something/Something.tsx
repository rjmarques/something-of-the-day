import React, { Component } from "react";
import { Button, Spinner } from "reactstrap";

import "./Something.scss";
import "../assets/img/github.svg";

import { ISomething } from "./SomethingAPI";

export interface IProps {
  getSomething: () => Promise<ISomething>;
}

export interface IState {
  loading: boolean;
  something?: ISomething;
}

export class Something extends Component<IProps, IState> {
  constructor(props: IProps) {
    super(props);

    this.state = {
      loading: false
    };
  }

  public componentDidMount() {
    this.updateSomething();
  }

  public render() {
    return (
      <div className="something-element">
        <div className="main-content">
          <header>
            <div className="title">Something of the day </div>
            <div>
              <img
                src={require("../assets/img/ascihappy.png")}
                alt="happy"
              ></img>
            </div>
          </header>
          <div className="something-box">
            <p>{this.state.something?.Text}</p>
            {this.state.loading ? (
              <Spinner type="grow" color="info" />
            ) : (
              <noscript />
            )}
          </div>
          <div className="button-area">
            <Button
              className="something-button"
              onClick={this.updateSomething}
              color="primary"
            >
              Ooh La La{" "}
              <span role="img" aria-label="kiss">
                💋
              </span>
            </Button>
          </div>
        </div>
        <footer>
          <div className="icons">
            <span>
              <a href="https://github.com/rjmarques" target="blank">
                <img src={require("../assets/img/github.svg")} alt="github" />
              </a>
            </span>
            <span>
              <a href="https://twitter.com/lgst_something" target="blank">
                <img src={require("../assets/img/twitter.svg")} alt="twitter" />
              </a>
            </span>
          </div>
          <div>
            <span>
              © Copyright {new Date().getFullYear()}{" "}
              <a href="https://ricardomarques.dev" target="blank">
                Ricardo Marques
              </a>
            </span>
          </div>
        </footer>
      </div>
    );
  }

  public updateSomething = async () => {
    try {
      if (this.state.loading) {
        return;
      }

      this.setState({ loading: true });

      const something = await this.props.getSomething();
      if (something) {
        this.setState({ something });
      }
    } catch (error) {
      console.log(error);
    } finally {
      this.setState({ loading: false });
    }
  };
}

export default Something;
