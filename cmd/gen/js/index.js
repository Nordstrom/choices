/** main gets the experiments
  * data format
  * []{
  *   namespaceName string
  *   experimentName string
  *   labels []string
  *   params map[string][]{
  *     name string
  *     value string
  *   }
  * }
  */
(function main() {
  const cookiePrefix = 'ExperimentId=';
  const start = document.cookie.indexOf(cookiePrefix);
  const end = document.cookie.indexOf(';', start);
  const id = document.cookie.slice(start+cookiePrefix.length, end);
  let req = new Request("{{.}}");
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
    json.forEach(experiment => {
      const expSection = document.createElement('section');
      expSection.style.border = '1px black solid';
      expSection.style.margin = '-1px 0px 0px -1px';
      const expTitle = document.createElement('h1');
      expTitle.textContent = experiment.experimentName;
      expSection.appendChild(expTitle);
      Object.keys(experiment.params).forEach(cookie => {
        const paramDiv = document.createElement('div');
        const experience = experiment.params[cookie].reduce((pv, param) => {
          pv.push(`${param.name}: ${param.value}`);
          return pv;
        }, []);
        const buttonText = experience.join(", ");
        const button = document.createElement('button');
        button.style.height = '40px';
        button.addEventListener('click', ev => {
          document.cookie = `experimentsdev=ExperimentId=${cookie};domain=.dev.nordstrom.com;path=/`;
          document.cookie = `experiments=ExperimentId=${cookie};domain=.nordstrom.com;path=/`;
          location.reload();
        });
        button.textContent = buttonText;
        paramDiv.appendChild(button);
        expSection.appendChild(paramDiv);
      });
      div.appendChild(expSection);
    });
    document.body.appendChild(div);
  });
}());
