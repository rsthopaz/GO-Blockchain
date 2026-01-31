package main
import "strings"
import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func main() {
  initDB()

  http.HandleFunc("/assets", assetHandler)
  http.HandleFunc("/assets/", assetByIDHandler)

  log.Println("Server running at :8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodPost:
    createAsset(w, r)
  case http.MethodGet:
    listAssets(w, r)
  default:
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  }
}

func createAsset(w http.ResponseWriter, r *http.Request) {
  var asset Asset
  json.NewDecoder(r.Body).Decode(&asset)

  asset.ID = uuid.New()
  asset.Status = "ACTIVE"

  query := `
    INSERT INTO assets 
    (id, asset_code, asset_name, category, value, status)
    VALUES ($1,$2,$3,$4,$5,$6)
  `

  _, err := db.Exec(
    query,
    asset.ID,
    asset.Code,
    asset.Name,
    asset.Category,
    asset.Value,
    asset.Status,
  )

  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  json.NewEncoder(w).Encode(asset)
}

func updateAsset(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
  var asset Asset
  json.NewDecoder(r.Body).Decode(&asset)

  query := `
    UPDATE assets
    SET asset_code=$1,
        asset_name=$2,
        category=$3,
        value=$4,
        updated_at=NOW()
    WHERE id=$5 AND status != 'DELETED'
  `

  res, err := db.Exec(
    query,
    asset.Code,
    asset.Name,
    asset.Category,
    asset.Value,
    id,
  )

  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  rows, _ := res.RowsAffected()
  if rows == 0 {
    http.Error(w, "Asset not found", 404)
    return
  }

  json.NewEncoder(w).Encode(map[string]string{
    "message": "Asset updated",
  })
}

func deleteAsset(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
  query := `
    UPDATE assets
    SET status='DELETED',
        updated_at=NOW()
    WHERE id=$1
  `

  res, err := db.Exec(query, id)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  rows, _ := res.RowsAffected()
  if rows == 0 {
    http.Error(w, "Asset not found", 404)
    return
  }

  json.NewEncoder(w).Encode(map[string]string{
    "message": "Asset deleted",
  })
}


func assetByIDHandler(w http.ResponseWriter, r *http.Request) {
  idStr := strings.TrimPrefix(r.URL.Path, "/assets/")
  id, err := uuid.Parse(idStr)
  if err != nil {
    http.Error(w, "Invalid asset ID", http.StatusBadRequest)
    return
  }

  switch r.Method {
  case http.MethodPut:
    updateAsset(w, r, id)
  case http.MethodDelete:
    deleteAsset(w, r, id)
  default:
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  }
}


func listAssets(w http.ResponseWriter, r *http.Request) {
  var assets []Asset
  query := `
  SELECT * FROM assets 
  WHERE status != 'DELETED'
  ORDER BY created_at DESC
`
  err := db.Select(&assets, query)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  json.NewEncoder(w).Encode(assets)
}
