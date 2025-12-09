# Pokedex Project - How It Works

## What This Project Does

This is a command-line Pokedex application that lets you explore Pokemon locations by connecting to the [PokeAPI](https://pokeapi.co). You can navigate forward and backward through pages of location data. Think of it like browsing a book - you can go to the next page or previous page of locations.

## Project Structure

```
pokedex/
├── main.go              # Entry point - starts the program
├── repl.go              # The command loop - reads your input
├── command_help.go      # Shows available commands
├── command_exit.go      # Exits the program
├── command_map.go       # Navigates location pages
├── pokecache.go         # (Currently empty - for caching)
└── settings/            # API-related code
    ├── client.go        # HTTP client setup
    ├── api.go           # Base API URL
    ├── loc_list.go      # Makes API calls
    └── loc_data.go      # Data structures for API responses
```

---

## How The Program Starts (main.go)

**File:** `main.go:9-16`

```go
func main() {
    pokeClient := pokeapi.NewClient(5 * time.Second)
    cfg := &config{
        pokeapiClient: pokeClient,
    }
    startRepl(cfg)
}
```

### What happens here:

1. **Creates an HTTP client** (`pokeClient`) with a 5-second timeout
   - The timeout means "give up if the API doesn't respond within 5 seconds"
   - This prevents your program from waiting forever

2. **Creates a config struct** that holds:
   - The API client (for making requests)
   - URLs for next/previous pages (stored later)

3. **Starts the REPL** (Read-Eval-Print Loop) - the interactive command loop

---

## The REPL - Your Command Loop (repl.go)

**File:** `repl.go:18-43`

The REPL is the heart of the program. It's a loop that:
1. Prints `Pokedex > ` prompt
2. Waits for your input
3. Executes the command
4. Repeats forever

### Step-by-Step Flow:

```go
func startRepl(cfg *config) {
    reader := bufio.NewScanner(os.Stdin)  // Creates input reader
    for {                                  // Infinite loop
        fmt.Print("Pokedex > ")           // Show prompt
        reader.Scan()                      // Wait for user input

        words := cleanInput(reader.Text()) // Convert to lowercase
        commandName := words[0]            // First word is command

        command, exists := getCommands()[commandName]
        if exists {
            err := command.callback(cfg)   // Run the command
            if err != nil {
                fmt.Println(err)          // Show errors
            }
        }
    }
}
```

### Available Commands:

The `getCommands()` function (repl.go:57-80) returns a map of commands:
- **help** - Shows all available commands
- **exit** - Exits the program
- **map** - Shows next page of locations
- **mapb** - Shows previous page of locations

---

## Understanding the HTTP Client (settings/client.go)

**File:** `settings/client.go:8-20`

```go
type Client struct {
    httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
    return Client{
        httpClient: http.Client{
            Timeout: timeout,
        },
    }
}
```

This creates a **wrapper** around Go's standard HTTP client. The wrapper:
- Stores an `http.Client` that will make all web requests
- Sets a timeout to prevent hanging requests
- Will be used by other functions to call the API

---

## How API Calls Work (settings/loc_list.go)

This is the most important part for understanding API calls!

**File:** `settings/loc_list.go:9-37`

```go
func (c *Client) ListLocations(pageURL *string) (RespShallowLocations, error) {
    // Step 1: Determine which URL to use
    url := baseURL + "/location-area"
    if pageURL != nil {
        url = *pageURL
    }

    // Step 2: Create HTTP GET request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return RespShallowLocations{}, err
    }

    // Step 3: Send the request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return RespShallowLocations{}, err
    }
    defer resp.Body.Close()  // Always close response body

    // Step 4: Read the response body
    dat, err := io.ReadAll(resp.Body)
    if err != nil {
        return RespShallowLocations{}, err
    }

    // Step 5: Parse JSON into Go struct
    locationResp := RespShallowLocations{}
    err = json.Unmarshal(dat, &locationResp)
    if err != nil {
        return RespShallowLocations{}, err
    }

    // Step 6: Return the parsed data
    return locationResp, nil
}
```

### Breaking Down The API Call:

#### Step 1: Build the URL
- Default URL: `https://pokeapi.co/api/v2/location-area`
- If you're navigating pages, use the URL from the previous response
- The `*string` means it can be `nil` (no value)

#### Step 2: Create the Request
- `http.NewRequest("GET", url, nil)` creates a GET request
- `"GET"` = the HTTP method (fetching data)
- `url` = where to send the request
- `nil` = no request body (GET requests don't send data)

#### Step 3: Send the Request
- `c.httpClient.Do(req)` actually sends the request over the internet
- This is where the program waits for PokeAPI to respond
- Returns a `response` object with status code, headers, and body

#### Step 4: Read the Response
- `io.ReadAll(resp.Body)` reads all bytes from the response
- The response body is the actual data (JSON) from the API
- `defer resp.Body.Close()` ensures we clean up when done

#### Step 5: Parse JSON
- `json.Unmarshal(dat, &locationResp)` converts JSON bytes to a Go struct
- Takes raw bytes (`dat`) and fills in the struct (`locationResp`)
- Go's struct tags (like `` `json:"count"` ``) tell it which JSON field maps to which struct field

#### Step 6: Return Data
- Returns the filled struct or an error if something went wrong
- Go functions can return multiple values (data and error)

---

## Understanding the Data Structure (settings/loc_data.go)

**File:** `settings/loc_data.go:3-11`

```go
type RespShallowLocations struct {
    Count    int     `json:"count"`
    Next     *string `json:"next"`
    Previous *string `json:"previous"`
    Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}
```

This struct matches the JSON response from PokeAPI:

```json
{
  "count": 1036,
  "next": "https://pokeapi.co/api/v2/location-area?offset=20&limit=20",
  "previous": null,
  "results": [
    {
      "name": "canalave-city-area",
      "url": "https://pokeapi.co/api/v2/location-area/1/"
    }
  ]
}
```

### Field Breakdown:

- **Count**: Total number of locations available
- **Next**: URL for the next page (nil if on last page)
- **Previous**: URL for previous page (nil if on first page)
- **Results**: Array of location objects with name and URL

The `*string` (pointer) allows the field to be `nil` when there's no next/previous page.

---

## How the Map Commands Work (command_map.go)

### Forward Navigation (map command)

**File:** `command_map.go:8-21`

```go
func commandMapf(cfg *config) error {
    // Make API call (nil means first page)
    locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationsURL)
    if err != nil {
        return err
    }

    // Save next/previous URLs for future navigation
    cfg.nextLocationsURL = locationsResp.Next
    cfg.prevLocationsURL = locationsResp.Previous

    // Print all location names
    for _, loc := range locationsResp.Results {
        fmt.Println(loc.Name)
    }
    return nil
}
```

### Backward Navigation (mapb command)

**File:** `command_map.go:23-40`

```go
func commandMapb(cfg *config) error {
    // Check if we're already on first page
    if cfg.prevLocationsURL == nil {
        return errors.New("you're on the first page")
    }

    // Make API call with previous page URL
    locationResp, err := cfg.pokeapiClient.ListLocations(cfg.prevLocationsURL)
    if err != nil {
        return err
    }

    // Update URLs
    cfg.nextLocationsURL = locationResp.Next
    cfg.prevLocationsURL = locationResp.Previous

    // Print locations
    for _, loc := range locationResp.Results {
        fmt.Println(loc.Name)
    }
    return nil
}
```

---

## Complete Flow: From Command to API to Display

Let's trace what happens when you type `map`:

1. **User types "map" in REPL** (repl.go:22)
   - Input is read and cleaned (converted to lowercase)

2. **Command is looked up** (repl.go:31)
   - `getCommands()["map"]` returns the map command struct

3. **Callback function runs** (repl.go:33)
   - `commandMapf(cfg)` is called

4. **API call is made** (command_map.go:9)
   - `ListLocations()` creates HTTP GET request
   - Request is sent to `https://pokeapi.co/api/v2/location-area`
   - Response comes back as JSON

5. **JSON is parsed** (loc_list.go:30-31)
   - Raw bytes converted to `RespShallowLocations` struct

6. **URLs are saved** (command_map.go:14-15)
   - Next and previous URLs stored in config for later use

7. **Results are displayed** (command_map.go:17-19)
   - Loop through each location and print name

8. **Back to REPL** (repl.go:20)
   - Prompt shows again, waiting for next command

---

## Key Go Concepts Used

### Pointers (`*string`)
- A pointer stores the memory address of a value
- `*string` can be `nil` (no value) or point to a string
- Used for Next/Previous URLs because they might not exist

### Error Handling
- Go functions return `(result, error)`
- Always check `if err != nil` after function calls
- Errors bubble up: if API call fails, return error to caller

### Structs
- Like objects in other languages
- Hold related data together
- `Client` struct holds the HTTP client
- `config` struct holds client and navigation state

### Methods
- `func (c *Client) ListLocations()` is a method on `Client`
- `c` is the receiver (like `this` or `self` in other languages)
- Lets you call `client.ListLocations()`

### JSON Tags
- `` `json:"count"` `` tells Go which JSON field maps to which struct field
- `json.Unmarshal` uses these tags to fill in data

### Defer
- `defer resp.Body.Close()` means "run this when function exits"
- Ensures cleanup happens even if errors occur

---

## What's Missing (pokecache.go)

The `pokecache.go` file is currently empty. It's likely meant to:
- Cache API responses to avoid repeated requests
- Store responses in memory with expiration times
- Speed up navigation when going back to previously visited pages

---

## How to Extend This Project

Here are ideas for learning more:

1. **Implement caching** - Store API responses so you don't re-fetch
2. **Add more commands** - Get details about a specific location
3. **Add error handling** - Better messages when API is down
4. **Add tests** - Write tests for your commands
5. **Fetch Pokemon data** - Expand beyond just locations

---

## Common API Call Pattern in Go

This project follows a standard pattern you'll see everywhere:

```go
// 1. Create request
req, err := http.NewRequest("GET", url, nil)
if err != nil {
    return err
}

// 2. Send request
resp, err := client.Do(req)
if err != nil {
    return err
}
defer resp.Body.Close()

// 3. Read response
data, err := io.ReadAll(resp.Body)
if err != nil {
    return err
}

// 4. Parse JSON
var result MyStruct
err = json.Unmarshal(data, &result)
if err != nil {
    return err
}

// 5. Use the data
return result, nil
```

This pattern works for any REST API!

---

## Questions to Help You Learn

As you explore this code, try to answer:

1. What happens if the API call takes more than 5 seconds?
2. Why do we use `*string` instead of `string` for Next/Previous?
3. What happens if you type an invalid command?
4. How would you add a command to show a specific location's details?
5. Why do we save the Next/Previous URLs in the config?

---

Happy coding! This project is a great foundation for learning Go and API interactions.
