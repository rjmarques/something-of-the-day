import React from "react";
import renderer from "react-test-renderer";
import Something from "./Something";
import { ISomething } from "./SomethingAPI";

it("renders correctly", async () => {
  const mockSomething: ISomething = {
    Id: 1,
    CreatedAt: "today",
    Text: "Let's put a smile on that face!"
  };
  const getSomething = jest
    .fn()
    .mockReturnValue(Promise.resolve(mockSomething));

  const component = renderer.create(<Something getSomething={getSomething} />);

  let tree = component.toJSON();
  expect(tree).toMatchSnapshot();

  await Promise.resolve();

  tree = component.toJSON();
  expect(tree).toMatchSnapshot();
});
