"use client";

import { FC, SyntheticEvent, useCallback, useMemo, useState } from "react";
import { Autocomplete, InputAdornment, OutlinedInput } from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";

// Query の候補を表示するための interface.
interface QueryInputOption {
  id?: string;
  label: string;
}

export const Search: FC = () => {
  const [value, setValue] = useState<string | QueryInputOption>("");

  const [options, setOptions] = useState<QueryInputOption[]>([]);

  // useMemo を使う必要がある。
  // see: https://mui.com/material-ui/react-autocomplete/#controlled-states
  const filteredOptions = useMemo(() => options, [options]);

  const onChange = useCallback(
    async (
      _: SyntheticEvent<Element, Event>,
      newValue: string | QueryInputOption | null
    ) => {
      console.log(`onChange: ${newValue}`);
      if (!newValue) {
        setOptions([]);
        return;
      }

      setValue(newValue);
      const q = typeof newValue === "string" ? newValue : newValue.label;
      if (newValue) {
        try {
          const response = await fetch(
            `http://localhost:8085/auto-complete?query=${q}`
          );
          if (!response.ok) {
            throw new Error("Failed to fetch data");
          }
          const data = await response.json();
          if (data.suggestions) {
            setOptions(data.suggestions);
          }
        } catch (error) {
          console.error("Error fetching data:", error);
        }
      }
    },
    []
  );

  return (
    <>
      <Autocomplete
        value={value}
        onChange={onChange}
        id="search-with-auto-complete"
        filterOptions={(x) => x}
        options={filteredOptions}
        autoComplete
        includeInputInList
        filterSelectedOptions
        noOptionsText={"候補なし"}
        getOptionLabel={(option: QueryInputOption | string) => {
          if (typeof option === "string") {
            return option;
          }
          return option.label;
        }}
        // console の warning
        // useAutocomplete.js:188 MUI: The value provided to Autocomplete is invalid. None of the options match with `"ta"`
        // を防ぐ。
        // see: https://stackoverflow.com/questions/61947941/material-ui-autocomplete-warning-the-value-provided-to-autocomplete-is-invalid#:~:text=Also%20when%20you%20want%20to%20build%20a%20searcher%20where%20value%20you%20write%20is%20not%20necesary%20the%20same%20as%20the%20options%20you%20can%20set%20freeSolo%20to%20true%20and%20the%20warning%20will%20disapear
        freeSolo
        selectOnFocus
        clearOnBlur
        handleHomeEndKeys
        // リストの表示をカスタマイズしたい時。
        renderOption={(props, option) => {
          const label = typeof option === "string" ? option : option.label;
          // keyをユニークにするためにdata-option-indexか他のユニークな属性を使用
          return (
            <li
              {...props}
              key={label}
              style={{
                ...props.style,
                color: "gray",
              }}
            >
              {label}
            </li>
          );
        }}
        renderInput={(params) => (
          <div ref={params.InputProps.ref}>
            <OutlinedInput
              fullWidth={true}
              rows={1}
              onChange={(event) => {
                onChange(event, event.target.value);
              }}
              startAdornment={
                <InputAdornment position="start">
                  <SearchIcon sx={{ fontSize: 20, mt: 0.25, ml: 0.5 }} />
                </InputAdornment>
              }
              inputProps={{
                ...params.inputProps,
              }}
              sx={{
                mt: 1,
                height: 48,
                fontWeight: "normal",
              }}
            />
          </div>
        )}
      />
    </>
  );
};

export default Search;
