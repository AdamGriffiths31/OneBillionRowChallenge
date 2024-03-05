#!/bin/bash
echo -ne "Start\r"
start_time=$(date +%s)

awk -F';' '{
    min[$1] = ($1 in min) ? (min[$1] < $2 ? min[$1] : $2) : $2;
    max[$1] = ($1 in max) ? (max[$1] > $2 ? max[$1] : $2) : $2;
    sum[$1] += $2;
    count[$1]++;
} END {
    printf("{");
    sep = "";
    PROCINFO["sorted_in"] = "@ind_str_asc"; # Sort by city name
    for (city in min) {
        mean = sum[city] / count[city];
        printf("%s%s=%s/%.1f/%s", sep, city, min[city], mean, max[city]);
        sep = ", ";
    }
    printf("}\n");
}' "$1" | sort  > output.txt

end_time=$(date +%s)
elapsed_time=$((end_time - start_time))

# Print elapsed time
echo "Elapsed time: $elapsed_time seconds"

read -p "Press Enter to exit"
