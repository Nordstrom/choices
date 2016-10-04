(function main() {
  const team = prompt("Enter Elwin team:");
  const cookiePrefix = 'ExperimentId=';
  const start = document.cookie.indexOf(cookiePrefix);
  const end = document.cookie.indexOf(';', start);
  const id = document.cookie.slice(start+cookiePrefix.length, end);
  let req = new Request(`http://elwin-test.ttoapps.aws.cloud.nordstrom.net/gen?label=${team}`);
  fetch(req)
  .then(resp => resp.json())
  .then(json => {
    const div = document.createElement('div');
    div.style.zIndex = '1';
    div.style.position = 'absolute';
    div.style.left = '10px';
    div.style.top = '10px';
    div.style.backgroundColor = 'rgba(255, 255, 255, 0.7)';
    div.style.padding = '10px';
    const data = Object.keys(json).reduce((pv, key) => {
      const params = key.split('.');
      pv[params[1]] = pv[params[1]] || {};
      pv[params[1]][params[2]] = pv[params[1]][params[2]] || {};
      pv[params[1]][params[2]][params[3]] = json[key];
      return pv
    }, {});
    Object.keys(data).forEach(experiment => {
      const expSection = document.createElement('section');
      expSection.style.border = '1px black solid';
      expSection.style.margin = '-1px 0px 0px -1px';
      const expTitle = document.createElement('h1');
      expTitle.textContent = experiment;
      expSection.appendChild(expTitle);
      Object.keys(data[experiment]).forEach(param => {
        const paramDiv = document.createElement('div');
        const paramTitle = document.createElement('h2');
        paramTitle.textContent = param;
        paramDiv.appendChild(paramTitle);
        Object.keys(data[experiment][param]).forEach(paramValue => {
          const paramValueDiv = document.createElement('div');
          const button = document.createElement('button');
          button.style.height = '40px';
          button.addEventListener('click', ev => {
            document.cookie = `experimentsdev=ExperimentId=${data[experiment][param][paramValue]};domain=.dev.nordstrom.com;path=/`;
            document.cookie = `experiments=ExperimentId=${data[experiment][param][paramValue]};domain=.nordstrom.com;path=/`;
            location.reload();
          });
          button.textContent = paramValue;
          paramValueDiv.appendChild(button);
          paramDiv.appendChild(paramValueDiv);
        });
        expSection.appendChild(paramDiv);
      });
      div.appendChild(expSection);
    });
    document.body.appendChild(div);
  });
}());
