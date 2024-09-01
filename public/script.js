const WORLDS_LIST = "/worlds.json";

/**
 * @param {string} worldName
 * @param {string} worldType
 * @returns {string}
 */
function composeLinkToWorldStatsTimeline(worldName, worldType) {
  return `/worlds/${worldType}/${worldName}.json`;
}

// const dateTimeFormat = new Intl.DateTimeFormat("pl-PL", {
//   day: "2-digit",
//   month: "short",
//   year: "numeric",
//   hour: "2-digit",
//   minute: "2-digit",
// });
// function formatDateTime(s) {
//   return dateTimeFormat.format(new Date(s * 1e3));
// }

/**
 * l: [timestamp, playersOnline][]
 * @typedef {{
 * l: [number, number][]
 * }} WorldStatsTimeline */

/** @typedef {{
 * name: string,
 * type: string,
 * }} World */

/** @typedef {{
 * worlds: World[]
 * }} WorldsList */

class App {
  /** @type HTMLSelectElement */
  #worldsList = document.querySelector("#worldsList");
  /** @type HTMLDivElement */
  #worldStatsContainer = document.querySelector("#worldStatsContainer");

  constructor() {
    console.log(this.#worldsList);
    console.log(this.#worldStatsContainer);
  }

  async setup() {
    await this.#setupWorldsList();
  }

  async #setupWorldsList() {
    const worldsList = await this.#fetchWorldsList();
    const worlds = worldsList.worlds.sort(
      (a, b) => b.type.localeCompare(a.type) || a.name.localeCompare(b.name),
    );
    this.#worldsList.innerHTML = "";
    for (const world of worlds) {
      const option = document.createElement("option");
      option.value = composeLinkToWorldStatsTimeline(world.name, world.type);
      option.innerText = world.name;
      this.#worldsList.appendChild(option);
    }

    if (worlds.length > 0) {
      this.#worldsList.selectedIndex = 0;
      const firstWorldLink = this.#worldsList.options[0].value;
      const worldStatsTimeline =
        await this.#fetchWorldStatsTimeline(firstWorldLink);
      this.#renderChart(worldStatsTimeline);
    }

    this.#worldsList.addEventListener("change", async (e) => {
      this.#worldStatsContainer.innerHTML = "";
      const link = e.target.value;
      const worldStatsTimeline = await this.#fetchWorldStatsTimeline(link);
      this.#renderChart(worldStatsTimeline);
    });
  }

  /**
   * @param {string} link
   * @returns {Promise<WorldStatsTimeline>}
   */
  async #fetchWorldStatsTimeline(link) {
    const res = await fetch(link);
    const data = await res.json();
    // TODO: Remove when the be is fixed
    // fix timezone
    const timeZoneOffset = 2 * 60 * 60 * 1000;
    // unix
    for (let i = 0; i < data.l.length; i++) {
      data.l[i][0] = data.l[i][0] * 1000 + timeZoneOffset;
    }
    return data;
  }

  /** @returns {Promise<WorldsList>} */
  async #fetchWorldsList() {
    const res = await fetch(WORLDS_LIST);
    const data = await res.json();
    return data;
  }

  /**
   * @param {WorldStatsTimeline} worldStatsTimeline
   */
  #renderChart(worldStatsTimeline) {
    this.#worldStatsContainer.innerHTML = "";

    Highcharts.stockChart(this.#worldStatsContainer, {
      rangeSelector: {
        selected: 1,
      },

      title: {
        text: "Postacie online",
      },

      series: [
        {
          name: "Postacie online",
          data: worldStatsTimeline.l,
        },
      ],
    });
  }
}

document.addEventListener("DOMContentLoaded", () => {
  const app = new App();
  app.setup();
  console.log("app setup finished");
});
