import { random_uint32 } from "./utils.js";

export const DIVISIONS = ["OPEN", "MODERN", "OLYMPIC", "TRADITIONAL"];
function rand_division() {
    return random_uint32() % DIVISIONS.length;
}

export const events = [
    {
        title: "Thurston vs Jesuit",
        leaders: new Map(),
        divisions: DIVISIONS,
        teams: {
            Thurston: {
                name: "Thurston",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Alex Morgan", division: rand_division(), score: 412, x_count: 6 },
                    { name: "Chris Nolan", division: rand_division(), score: 398, x_count: 3, },
                    { name: "Taylor Reed", division: rand_division(), score: 430, x_count: 11, },
                    { name: "Jordan Kim", division: rand_division(), score: 445, x_count: 9 },
                    { name: "Sam Ortega", division: rand_division(), score: 376, x_count: 2, },
                    { name: "Riley Chen", division: rand_division(), score: 459, x_count: 14, },
                ],
            },
            Jesuit: {
                name: "Jesuit",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Evan Brooks", division: rand_division(), score: 438, x_count: 10 },
                    { name: "Noah Patel", division: rand_division(), score: 471, x_count: 18, },
                    { name: "Liam O'Connor", division: rand_division(), score: 401, x_count: 5, },
                    { name: "Diego Ramirez", division: rand_division(), score: 455, x_count: 12, },
                    { name: "Mason Wright", division: rand_division(), score: 429, x_count: 7 },
                ],
            },
        },
    },

    {
        title: "Oregon Outdoor State Championship",
        leaders: new Map(),
        divisions: DIVISIONS,
        teams: {
            "Central Catholic": {
                name: "Central Catholic",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Sophia Lin", division: rand_division(), score: 472, x_count: 19 },
                    { name: "Marcus Hill", division: rand_division(), score: 488, x_count: 22, },
                    { name: "Tyler Bennett", division: rand_division(), score: 405, x_count: 6, },
                    { name: "Avery Collins", division: rand_division(), score: 451, x_count: 13, },
                    { name: "Julian Perez", division: rand_division(), score: 479, x_count: 20, },
                    { name: "Harper Nguyen", division: rand_division(), score: 392, x_count: 3, },
                    { name: "Kevin Foster", division: rand_division(), score: 436, x_count: 9 },
                ],
            },
            Jesuit: {
                name: "Jesuit",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Isabella Cruz", division: rand_division(), score: 469, x_count: 18, },
                    { name: "Andrew Park", division: rand_division(), score: 492, x_count: 24, },
                    { name: "Sean Murphy", division: rand_division(), score: 410, x_count: 7, },
                    { name: "Daniel Kim", division: rand_division(), score: 447, x_count: 11 },
                    { name: "Victor Gomez", division: rand_division(), score: 475, x_count: 19, },
                ],
            },
            "West Linn": {
                name: "West Linn",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Olivia Hart", division: rand_division(), score: 455, x_count: 14 },
                    { name: "Ryan Cooper", division: rand_division(), score: 399, x_count: 4, },
                    { name: "Ethan Zhou", division: rand_division(), score: 483, x_count: 21, },
                    { name: "Lucas Grant", division: rand_division(), score: 442, x_count: 10 },
                    { name: "Nina Petrova", division: rand_division(), score: 468, x_count: 17, },
                    { name: "Cole Harrison", division: rand_division(), score: 387, x_count: 2, },
                ],
            },
        },
    },

    {
        title: "OHSAL Public Open",
        leaders: new Map(),
        divisions: DIVISIONS,
        teams: {
            Springfield: {
                name: "Springfield",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Parker Young", division: rand_division(), score: 432, x_count: 8 },
                    { name: "Logan Reed", division: rand_division(), score: 467, x_count: 16, },
                    { name: "Dylan Moore", division: rand_division(), score: 390, x_count: 3, },
                    { name: "Elliot Baker", division: rand_division(), score: 446, x_count: 12, },
                    { name: "Miguel Santos", division: rand_division(), score: 458, x_count: 14, },
                ],
            },
            "North Salem": {
                name: "North Salem",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Adrian Fox", division: rand_division(), score: 419, x_count: 7 },
                    { name: "Ben Wallace", division: rand_division(), score: 401, x_count: 5, },
                    { name: "Ivan Petrov", division: rand_division(), score: 471, x_count: 18, },
                    { name: "Theo Ramirez", division: rand_division(), score: 437, x_count: 10, },
                    { name: "Kira Volkov", division: rand_division(), score: 460, x_count: 15, },
                    { name: "Jesse Long", division: rand_division(), score: 382, x_count: 2, },
                ],
            },
            McKenzie: {
                name: "McKenzie",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Carter Mills", division: rand_division(), score: 428, x_count: 9 },
                    { name: "Wyatt Dean", division: rand_division(), score: 395, x_count: 4, },
                    { name: "Hugo Laurent", division: rand_division(), score: 463, x_count: 17, },
                    { name: "Tariq Hassan", division: rand_division(), score: 441, x_count: 11, },
                    { name: "Miles Carter", division: rand_division(), score: 452, x_count: 13, },
                ],
            },
            Ashland: {
                name: "Ashland",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Noelle Wright", division: rand_division(), score: 435, x_count: 9, },
                    { name: "Grant Olson", division: rand_division(), score: 388, x_count: 3, },
                    { name: "Soren Dahl", division: rand_division(), score: 469, x_count: 18, },
                    { name: "Priya Nair", division: rand_division(), score: 449, x_count: 12 },
                    { name: "Leo Martinez", division: rand_division(), score: 457, x_count: 14, },
                ],
            },
        },
    },

    {
        title: "4A State Qualifiers",
        leaders: new Map(),
        divisions: DIVISIONS,
        teams: {
            Roseburg: {
                name: "Roseburg",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Hunter Blake", division: rand_division(), score: 423, x_count: 7 },
                    { name: "Spencer King", division: rand_division(), score: 392, x_count: 4, },
                    { name: "Dominic Ruiz", division: rand_division(), score: 468, x_count: 17, },
                    { name: "Evan Price", division: rand_division(), score: 441, x_count: 10 },
                    { name: "Jonah Klein", division: rand_division(), score: 459, x_count: 15, },
                ],
            },
            Silverton: {
                name: "Silverton",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Aiden Ross", division: rand_division(), score: 436, x_count: 11 },
                    { name: "Calvin Brooks", division: rand_division(), score: 401, x_count: 6, },
                    { name: "Mateo Alvarez", division: rand_division(), score: 472, x_count: 19, },
                    { name: "Toby Chen", division: rand_division(), score: 448, x_count: 12 },
                    { name: "Roman Novak", division: rand_division(), score: 461, x_count: 16, },
                    { name: "Finn Carter", division: rand_division(), score: 387, x_count: 2, },
                ],
            },
            Stayton: {
                name: "Stayton",
                score: 0,
                x_count: 0,
                members: [
                    { name: "Levi Turner", division: rand_division(), score: 429, x_count: 8 },
                    { name: "Brady Owens", division: rand_division(), score: 396, x_count: 4, },
                    { name: "Rafael Costa", division: rand_division(), score: 465, x_count: 17, },
                    { name: "Nolan Pierce", division: rand_division(), score: 444, x_count: 11, },
                    { name: "Emil Johansson", division: rand_division(), score: 458, x_count: 14, },
                ],
            },
        },
    },
];
