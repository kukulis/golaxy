class ApiClient {

    /**
     * @param id
     * @returns {Promise<Battle>}
     */
    async getBattle(id) {
        const response = await fetch('/api/battle');

        if (!response.ok) {
            throw new Error(`Failed to fetch battle: ${response.statusText}`);
        }

        const data = await response.json();

        return (new Battle()).updateFromDTO(data);
    }
}